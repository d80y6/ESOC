from opensearchpy import OpenSearch
import logging

logger = logging.getLogger(__name__)

class RAGManager:
    def __init__(self, hosts=["http://opensearch:9200"]):
        self.client = OpenSearch(
            hosts=hosts,
            http_compress=True,
            use_ssl=False,
            verify_certs=False,
            ssl_assert_hostname=False,
            ssl_show_warn=False,
        )

    def retrieve_context(self, query, index="normalized-events-*", limit=5):
        logger.info(f"Retrieving context for query: {query}")

        # Simple keyword search for demonstration.
        # In production, this would use vector embeddings and k-NN search.
        body = {
            "query": {
                "multi_match": {
                    "query": query,
                    "fields": ["message", "event.action", "rule_name"]
                }
            },
            "size": limit
        }

        try:
            response = self.client.search(index=index, body=body)
            hits = response['hits']['hits']
            return [hit['_source'] for hit in hits]
        except Exception as e:
            logger.error(f"Failed to retrieve from OpenSearch: {e}")
            return []

    def format_context(self, events):
        if not events:
            return "No relevant historical events found."

        formatted = "Relevant historical events:\n"
        for i, event in enumerate(events):
            formatted += f"{i+1}. Rule: {event.get('rule_name', 'Unknown')}, Time: {event.get('@timestamp')}, Host: {event.get('host', {}).get('hostname', 'N/A')}\n"
        return formatted
