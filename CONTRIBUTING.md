# How to contribute to Fireblaze Vault

When contributing to this project, please first discuss the change you wish to make via issue, email, or any other method with the maintainers before making a change. Please note we have a [code of conduct](./CODE_OF_CONDUCT.md), follow it in all your interactions with the project.

## Quick start

`$ make install` will go get dependencies for protocol buffers, grpc and ent. 

`$ make generate` will use `protoc` and `entc` to generate the api, swagger and db schema.

Protobufs are in the `proto/` folder. If you update the .proto files, regenerate using `make proto`

Database schemas are in the `ent/schema/` folder. Regenerate using `make schema`

Please don't commit generated code. Maintain .gitignore accordingly.

## Folder structure
`proto/` contains the protobuf specification and generated code.

`ent/` contains db `schema/` specification, the generated code lives in multiple locations inside the `ent/` folder.

`internal/` contains business logic inherent only to this project.  

`pkg/` contains reusable packages with tests. If you don't have tests yet, keep the package in `internal/`.

`cmd/` contains CLI applications.

`hack/` contains automation/convenience scripts to ease repetitive tasks.

`web/` contains UI implementations.

`kube/` contains the infrastructure code, mostly yaml manifests.

## Pull Request Process

1. Fork
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Add new feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create new Pull/Merge Request

## Sign Your Work

The sign-off is a simple line at the end of the explanation for a commit. All
commits needs to be signed. Your signature certifies that you wrote the patch or
otherwise have the right to contribute the material. The rules are pretty simple,
if you can certify points (a),(b),(c),(d) from ([developercertificate.org](http://developercertificate.org/))

Then you just add a line to every git commit message:

    Signed-off-by: Joe Smith <joe.smith@example.com>

If you set your `user.name` and `user.email` git configs, you can sign your
commit automatically with `git commit -s`.

Note: If your git config information is set properly then viewing the
`git log` information for your commit will look something like this:

```
Author: Joe Smith <joe.smith@example.com>
Date:   Thu Feb 2 11:41:15 2018 -0800

    Update README

    Signed-off-by: Joe Smith <joe.smith@example.com>
```

Notice the `Author` and `Signed-off-by` lines match. If they don't
your PR will be rejected by the automated DCO check.

## Issues

Issues are used as the primary method for tracking anything to do with the project.

### Issue Types

There are 4 types of issues (each with their own corresponding [label](#labels)):

- Question: These are support or functionality inquiries that we want to have a record of for
  future reference. Generally these are questions that are too complex or large to store in the
  Slack channel or have particular interest to the community as a whole. Depending on the discussion,
  these can turn into "Feature" or "Bug" issues.
- Proposal: Used for items (like this one) that propose new ideas or functionality that require a larger community discussion. This allows for feedback from others in the community before a
  feature is actually developed. This is not needed for small additions. Final word on whether or
  not a feature needs a proposal is up to the core maintainers. All issues that are proposals should
  both have a label and an issue title of "Proposal: [the rest of the title]." A proposal can become
  a "Feature" and does not require a milestone.
- Features: These track specific feature requests and ideas until they are complete. They can evolve
  from a "Proposal" or can be submitted individually depending on the size.
- Bugs: These track bugs with the code or problems with the documentation (i.e. missing or incomplete)

### Issue Lifecycle

The issue lifecycle is mainly driven by the core maintainers, but is good information for those
contributing to Fireblaze Vault. All issue types follow the same general lifecycle. Differences are noted below.

1. Issue creation
2. Triage
   - The maintainer in charge of triageing will apply the proper labels for the issue. This
     includes labels for priority, type, and metadata (such as "starter"). The only issue
     priority we will be tracking is whether or not the issue is "critical." If additional
     levels are needed in the future, we will add them.
   - (If needed) Clean up the title to succinctly and clearly state the issue. Also ensure
     that proposals are prefaced with "Proposal".
   - Add the issue to the correct milestone. If any questions come up, don't worry about
     adding the issue to a milestone until the questions are answered.
   - We attempt to do this process at least once per work day.
3. Discussion
   - "Feature" and "Bug" issues should be connected to the PR that resolves it.
   - Whoever is working on a "Feature" or "Bug" issue (whether a maintainer or someone from
     the community), should either assign the issue to them self or make a comment in the issue
     saying that they are taking it.
   - "Proposal" and "Question" issues should stay open until resolved or if they have not been
     active for more than 30 days. This will help keep the issue queue to a manageable size and
     reduce noise. Should the issue need to stay open, the `keep open` label can be added.
4. Issue closure