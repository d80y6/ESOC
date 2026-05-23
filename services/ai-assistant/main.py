from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
from typing import List, Optional, Dict, Any
from core.rag import RAGManager
import uvicorn
import os

app = FastAPI(title="OmniGuard AI Assistant")
rag = RAGManager()

class TriageRequest(BaseModel):
    alert_id: str
    alert_details: Dict[str, Any]

class TriageResponse(BaseModel):
    summary: str
    severity_score: int
    recommendations: List[str]
    context_found: bool

@app.get("/health")
def health():
    return {"status": "ok"}

@app.post("/triage", response_model=TriageResponse)
def triage_alert(req: TriageRequest):
    # 1. Retrieve context
    context_events = rag.retrieve_context(req.alert_details.get("rule_name", ""))
    context_str = rag.format_context(context_events)

    # 2. Orchestrate LLM (Mocking for this platform build)
    # In production, we'd call an OpenAI-compatible API or a local LLM (Ollama/vLLM)

    summary = f"Analysis of Alert {req.alert_id}: This detection was triggered by {req.alert_details.get('rule_name')}. "
    summary += "Cross-referencing with historical data, similar patterns were observed on other hosts. "
    summary += "The activity appears suspicious and warrants immediate investigation."

    return TriageResponse(
        summary=summary,
        severity_score=8,
        recommendations=[
            "Isolate the host immediately",
            "Reset the credentials for the involved user",
            "Check for established C2 connections"
        ],
        context_found=len(context_events) > 0
    )

if __name__ == "__main__":
    uvicorn.run(app, host="0.0.0.0", port=8088)
