from fastapi import FastAPI, BackgroundTasks, Request
from pydantic import BaseModel
from typing import List, Optional, Dict, Any
from core.executor import WorkflowExecutor
from core.auth import AuthMiddleware
import uvicorn

app = FastAPI(title="OmniGuard SOAR")
executor = WorkflowExecutor()

# Apply Auth Middleware
app.middleware("http")(AuthMiddleware(
    issuer_url="http://keycloak:8080/realms/omniguard",
    audience="omniguard-backend"
))

class Workflow(BaseModel):
    id: str
    name: str
    description: Optional[str] = None
    steps: List[Dict[str, Any]]

class TriggerRequest(BaseModel):
    workflow_id: str
    context: Dict[str, Any]

# In-memory store for demo
workflows = {}

@app.get("/health")
def health():
    return {"status": "ok"}

@app.post("/workflows")
def create_workflow(workflow: Workflow):
    workflows[workflow.id] = workflow.dict()
    return {"id": workflow.id, "status": "created"}

@app.post("/trigger")
def trigger_workflow(req: TriggerRequest, background_tasks: BackgroundTasks):
    workflow = workflows.get(req.workflow_id)
    if not workflow:
        return {"error": "Workflow not found"}, 404

    background_tasks.add_task(executor.execute, workflow, req.context)
    return {"status": "triggered"}

if __name__ == "__main__":
    uvicorn.run(app, host="0.0.0.0", port=8087)
