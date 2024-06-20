import datetime
import json
import os
import re

import git
from github import Github, Auth

import generate

ctx = json.loads(os.environ.get("GITHUB_CONTEXT"))
gh = Github(auth=Auth.Token(os.environ.get("GITHUB_TOKEN")))

# closed: do publish
if ctx["event"]["action"] == "closed":
    repo = gh.get_repo("hidetatz/hidetatz.github.io")
    issue = repo.get_issue(ctx["event"]["issue"]["number"])

    ts = issue.closed_at
    pattern = "20\d{2}年\d{1,2}月\d{1,2}日"
    result = re.search(pattern, issue.title)
    if result:
        ts = datetime.datetime.strptime(result.group(), "%Y年%m月%d日")
    generate.new_diary(issue.title, ts, issue.body)
    generate.generate()
    gh.close()

    repo = git.Repo()
    repo.git.commit(".", "-m", "update diary")
    origin = repo.remote(name="origin")
    origin.push()

# # edited: if the issue is already closed, do publish
# if ctx["event"]["action"] == "edited":
#     # publish
#     gh.close()
#     return
# 
# # reopened, deleted: unpublish
# if ctx["event"]["action"] in ["reopened", "deleted"]:
#     # unpublish
#     gh.close()
#     return

else:
    print(f"event {ctx['event']['action']} is ok not to handle")
    gh.close()
