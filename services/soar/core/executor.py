import requests
import logging

logger = logging.getLogger(__name__)

class WorkflowExecutor:
    def __init__(self):
        self.actions = {
            "http_request": self.http_request,
            "slack_notify": self.slack_notify,
        }

    def execute(self, workflow_data, context):
        logger.info(f"Executing workflow: {workflow_data['name']}")
        for step in workflow_data['steps']:
            action_type = step.get("type")
            action_func = self.actions.get(action_type)
            if action_func:
                try:
                    action_func(step.get("params"), context)
                except Exception as e:
                    logger.error(f"Step {step['name']} failed: {e}")
                    if step.get("on_failure") == "stop":
                        break
            else:
                logger.warning(f"Unknown action type: {action_type}")

    def http_request(self, params, context):
        method = params.get("method", "GET")
        url_template = params.get("url")

        # SSRF Protection: URL Allowlist (Demo)
        allowed_domains = ["slack.com", "hooks.slack.com", "api.pagerduty.com", "outlook.office.com"]

        # Safely render the URL template to avoid object attribute access
        try:
            import re
            def safe_sub(match):
                key = match.group(1)
                return str(context.get(key, match.group(0)))

            url = re.sub(r'\{([a-zA-Z0-9_-]+)\}', safe_sub, url_template)
        except Exception as e:
            logger.error(f"Safe template rendering failed: {e}")
            return

        from urllib.parse import urlparse
        parsed_url = urlparse(url)
        if parsed_url.netloc not in allowed_domains:
            logger.error(f"Access to domain {parsed_url.netloc} is forbidden (SSRF Protection)")
            return

        resp = requests.request(method, url, json=params.get("body"), timeout=10)
        logger.info(f"HTTP {method} to {url} returned {resp.status_code}")

    def slack_notify(self, params, context):
        webhook_url = params.get("webhook_url")
        message = params.get("message", "Alert triggered").format(**context)
        requests.post(webhook_url, json={"text": message})
        logger.info("Sent Slack notification")
