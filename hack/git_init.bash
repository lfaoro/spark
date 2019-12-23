rm -rf .git
git init
git remote add origin git@github.com:lfaoro/spark.git
git remote -v
git add --all
git commit -am "first commit"
git push --force origin master