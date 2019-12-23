// Copyright (c) 2019 Leonardo Faoro. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package vault

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"

	"github.com/lfaoro/spark/ent/card"
	"github.com/lfaoro/spark/internal"
	"github.com/lfaoro/spark/pkg/valid"
	pb "github.com/lfaoro/spark/proto/api/vault"
)

func (s Server) PutCard(ctx context.Context, r *pb.PutCardRequest) (w *pb.PutCardResponse, err error) {
	var ciphertext []byte

	// sanitize payment card number
	r.Number = strings.ReplaceAll(r.Number, " ", "")

	ok := valid.IsCreditCard(r.Number)
	if !ok {
		return nil,
			status.Errorf(codes.FailedPrecondition,
				"validation failed: %s", "not a valid credit card")
	}

	first6, last4 := getCardPrefixSuffix(r.Number)
	w = &pb.PutCardResponse{
		ExpiresOn: cardExpiresOn(r),
		FirstSix:  uint32(first6),
		LastFour:  uint32(last4),
		RequestIp: ipFrom(ctx),
		UserAgent: userAgentsFrom(ctx),
	}

	for service, ok := range s.enabledServices(ctx) {
		if !ok {
			continue
		}

		switch service {
		case pb.SERVICE_TOKENIZATION:
			w.Token = newToken()
			w.Hash = hashCard(r)
			w.AutoDeleteOn = autoDeleteOn(r)
			ciphertext, err = s.encryptCard(r)
			if err != nil {
				return nil, status.Errorf(codes.Internal,
					"failed to encrypt card data: %s", err)
			}
		// s.billing.Inc(pID, billing.Card)

		case pb.SERVICE_METADATA:
			w.Metadata = s.iinLookup(r.Number)
		// s.billing.Inc(billing.Metadata)
		case pb.SERVICE_RISK_CHECK:
			w.Risk = calculateRisk(r)
		// s.billing.Inc(billing.Risk)

		case pb.SERVICE_MPI_CHECK:
			w.Mpi = getMPI(r)
			// s.billing.Inc(billing.MPI)
		}
	}

	err = s.storeCard(ctx, w, ciphertext) // todo: don't store here
	if err != nil {
		code := codes.Internal
		if strings.Contains(err.Error(), "unique constraint failed") {
			err = errors.New("a card with these details has already been stored")
			code = codes.Aborted
		}
		_err := status.Errorf(code,
			"storage error: %s", err)
		return nil, _err
	}

	metricStoredCards.Inc()

	return w, nil
}

func (s Server) GetCard(ctx context.Context, r *pb.GetCardRequest) (w *pb.GetCardResponse, err error) {
	query, err := s.findCard(ctx, r.Token)
	if err != nil {
		log.Printf("failed query with error (%v) for token %s\n", err, r.Token)
		_err := status.Errorf(codes.NotFound,
			"failed query for token %v", r.Token)
		return nil, _err

	}
	w, err = s.decryptCard(query.CipheredCard)
	if err != nil {
		return nil,
			status.Errorf(codes.Internal,
				"decryption failed: %s", err)
	}

	tp, _ := ptypes.TimestampProto(query.AutoDeleteOn)
	w.AutoDeleteOn = tp

	metricQueryCard.Inc()

	return w, err
}

func (s Server) GetCardMetadata(ctx context.Context, r *pb.GetMetadataRequest) (*pb.CardMetadata, error) {
	q, err := s.findCardMetadata(ctx, r.Token)
	if err != nil {
		return nil,
			status.Errorf(codes.NotFound,
				"failed query: %s", err)
	}

	w := &pb.CardMetadata{
		Scheme:    q.Scheme,
		Brand:     q.Brand,
		Type:      q.Type,
		Currency:  q.Currency,
		IsPrepaid: q.IsPrepaid,
		Issuer: &pb.CardMetadata_Issuer{
			Name:        q.IssuerName,
			Url:         q.IssuerURL,
			Phone:       q.IssuerPhone,
			City:        q.IssuerCity,
			Country:     q.IssuerCountry,
			CountryCode: q.IssuerCountrycode,
			Latitude:    float32(q.IssuerLatitude),
			Longitude:   float32(q.IssuerLongitude),
			Map:         q.IssuerMap,
		},
	}

	metricQueryCard.Inc()

	return w, err
}

// DelCard removes all the data associated with a token from the database.
func (s Server) DelCard(ctx context.Context, r *pb.DelCardRequest) (*empty.Empty, error) {
	n, err := s.db.Card.
		Delete().
		Where(card.Token(r.Token)).
		Exec(ctx)
	if err != nil {
		return nil,
			status.Errorf(codes.NotFound,
				"failed query: %s", err)
	}

	if n == 0 {
		return nil, status.Errorf(codes.NotFound, "failed to find token: %v", r.Token)
	}

	metricDeleteCard.Inc()

	return &empty.Empty{}, err
}

func (s Server) iinLookup(number string) *pb.CardMetadata {
	if s.iin == nil {
		return &pb.CardMetadata{ServiceMessage: "metadata lookup: disabled"}
	}

	meta, err := s.iin.Lookup(number)
	if err != nil {
		s.err = fmt.Errorf("lookup: %w", err)
		grpclog.Error(err)
		return &pb.CardMetadata{
			ServiceMessage: fmt.Sprintf(
				"metadata lookup failed: %v", err),
		}
	}

	return &pb.CardMetadata{
		Scheme:    meta.CardScheme,
		Brand:     meta.CardBrand,
		Type:      meta.CardType,
		Currency:  meta.CardCurrency,
		IsPrepaid: meta.CardIsPrepaid,
		Issuer: &pb.CardMetadata_Issuer{
			Name:        meta.Issuer,
			Url:         meta.IssuerURL,
			Phone:       meta.IssuerPhone,
			City:        meta.IssuerCity,
			Country:     meta.IssuerCountry,
			CountryCode: meta.IssuerCountryCode,
			Latitude:    meta.IssuerLatitude,
			Longitude:   meta.IssuerLongitude,
			Map:         meta.IssuerMap,
		},
	}
}

// todo: implement
func getMPI(r *pb.PutCardRequest) *pb.MPI {
	return &pb.MPI{
		Enrolled: true,
		Eci:      2,
		Acs:      "DEMO: https://secure5.arcot.com/acspage/cap?RID=35325&VAA=B",
		Par:      "DEMO: eNpdU8tymzAU3ecrvMumYz1AgD2yZnDsTpMZx06TRdudLK5skhiIgNjust/Tr+qXVMIm4DDDoHvuOdKZexB/2hqA2SOo2oC4Ggz4AspSbmCQJpPrSELAZKApSOyFkWYaY+brJAq0pzTQa6ewmlX8Hd4aRZDoKMIs8WXiK2/NvIAFDPwwoRrCkV6fFVbzDqZM80yQIR5SjtqybRcyE4xFPmUhoZiGhPkRRw5tGQswaiuzqgUsJNXb9PZeMBpYFUfnsuvvwNzOBKa4fTg6QR0lkzsQ84PcFa/AUVN1TZXXWWWOgmLrpS26dm1exbaqijFC+/1+CKddhrnZII5cs7WOPnvnq9oBZf+wQ5qIxSzed+/8uPj9QO6flb98iiccOUbHT2QF1hlhmODRgARj6o0Z46jBezPaOd/i35+/hH7x7ATOQMconJf4hBLqKH2kN43aGMjUUYxCN4626ghwKPIMrMbm+7HuGYZSCevPfT4m83kQ/ObbRcCqsnHls7V++fVAHoPlnb/9eVMuVz+m8fzrNN5MXOwN6cJHaoMiETkZSbvUOGr3t0e7v7i5A+h8CcQVR5cX5D/FxN2u",
	}
}

func getCardPrefixSuffix(number string) (int, int) {
	first6, err := strconv.Atoi(number[:6])
	last4, err := strconv.Atoi(number[len(number)-4:])
	if err != nil {
		return 0, 0
	}
	return first6, last4
}

func cardExpiresOn(r *pb.PutCardRequest) *timestamp.Timestamp {
	t := time.Date(int(r.ExpYear), time.Month(r.ExpMonth)+1, 1, 0, 0, 0, 1, time.UTC)
	expTime, err := ptypes.TimestampProto(t)
	log.LogIfErr(err)
	return expTime
}

func userAgentsFrom(ctx context.Context) string {
	md, _ := metadata.FromIncomingContext(ctx)
	return strings.Join(md.Get("grpcgateway-user-agent"), " ")
}

func ipFrom(ctx context.Context) string {
	p, ok := peer.FromContext(ctx)
	if !ok {
		return ""
	}
	ip, _, err := net.SplitHostPort(p.Addr.String())
	if err != nil {
		return ""
	}
	return ip
}

func autoDeleteOn(card *pb.PutCardRequest) *timestamp.Timestamp {
	t := time.Date(2100, 0, 0, 0, 0, 0, 0, time.UTC)
	tp, _ := ptypes.TimestampProto(t)
	switch card.AutoDelete {
	case pb.PutCardRequest_NEVER:
		return tp

	case pb.PutCardRequest_EXPIRY_DATE:
		return cardExpiresOn(card)

	case pb.PutCardRequest_THREE_MONTHS:
		t := time.Now().AddDate(0, 3, 0)
		tp, _ := ptypes.TimestampProto(t)
		return tp

	case pb.PutCardRequest_SIX_MONTHS:
		t := time.Now().AddDate(0, 6, 0)
		tp, _ := ptypes.TimestampProto(t)
		return tp

	case pb.PutCardRequest_TWELVE_MONTHS:
		t := time.Now().AddDate(0, 12, 0)
		tp, _ := ptypes.TimestampProto(t)
		return tp
	}
	return tp
}

// todo: implement
func calculateRisk(card *pb.PutCardRequest) *pb.Risk {
	// todo: implement
	// if risk > 70 {
	// 	metricHighRisk.Inc()
	// }

	return &pb.Risk{
		Score: 33.27,
		Triggers: []pb.Risk_Trigger{
			pb.Risk_BLACKLIST,
			pb.Risk_FINGERPRINT,
			pb.Risk_GEOLOCATION,
		},
	}
}

func hashCard(r *pb.PutCardRequest) string {
	var buf strings.Builder
	buf.WriteString(r.Number)
	buf.WriteString(string(r.ExpMonth))
	buf.WriteString(string(r.ExpYear))
	cc := buf.String()
	hash := sha256.New()
	hash.Write([]byte(cc))
	sum := hash.Sum(nil)

	return base64.StdEncoding.EncodeToString(sum)
}

func newToken() string {
	return "tok_" + uuid.New().String()
}

func (s Server) enabledServices(ctx context.Context) map[pb.SERVICE]bool {
	services := make(map[pb.SERVICE]bool)

	id := projectIDFrom(ctx)
	if id == 0 {
		return services
	}

	project, err := s.db.Project.Get(ctx, id)
	if err != nil {
		log.Println(err)
		return services
	}

	services[pb.SERVICE_TOKENIZATION] = project.Tokenization
	services[pb.SERVICE_METADATA] = project.Metadata
	services[pb.SERVICE_RISK_CHECK] = project.Risk
	services[pb.SERVICE_MPI_CHECK] = project.Mpi

	return services
}

func projectIDFrom(ctx context.Context) int {
	id, ok := internal.GetValueFrom(ctx, "project_id").(int)
	if !ok {
		log.Println("failed to get project_id from context")
		return 0
	}
	return id
}
