import os

ctx = json.loads(os.environ.get("GITHUB_CONTEXT"))

# closed: do publish
if ctx["event"]["action"] == "closed":
    # publish
    return

# edited: if the issue is already closed, do publish
if ctx["event"]["action"] == "edited":
    # publish
    return

# reopened, deleted: unpublish
if ctx["event"]["action"] in ["reopened", "deleted"]:
    # unpublish
    return

print(f"event {ctx['event']['action']} is ok not to handle")
