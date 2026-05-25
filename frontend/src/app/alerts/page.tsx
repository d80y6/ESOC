"use client";

import { useState, useEffect } from "react";
import { Search, Filter, MoreHorizontal, RefreshCw } from "lucide-react";

export default function AlertConsole() {
  const [searchTerm, setSearchTerm] = useState("");
  const [alerts, setAlerts] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);

  const fetchAlerts = async () => {
    setLoading(true);
    try {
      // In a real ISP environment, this would call the Alerting Service API
      // For this implementation, we fetch from the Search service using the 'alerts' filter
      const response = await fetch("/api/search", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          query: "rule_name: *",
          size: 20
        }),
      });
      const data = await response.json();
      setAlerts(data.hits?.hits?.map((h: any) => ({
        id: h._id,
        time: h._source["@timestamp"],
        rule: h._source["rule_name"],
        severity: h._source["severity"] || "Medium",
        host: h._source["host.name"] || "N/A",
        status: h._source["status"] || "Open"
      })) || []);
    } catch (err) {
      console.error("Failed to fetch alerts", err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchAlerts();
  }, []);

  return (
    <div className="min-h-screen bg-slate-950 text-slate-50 p-8">
      <header className="mb-8 flex justify-between items-center">
        <div>
          <h1 className="text-3xl font-bold">Alert Console</h1>
          <p className="text-slate-400">Review and triage security detections</p>
        </div>
        <div className="flex gap-4">
           <div className="relative">
              <Search className="absolute left-3 top-2.5 h-4 w-4 text-slate-500" />
              <input
                type="text"
                placeholder="Search alerts..."
                className="bg-slate-900 border border-slate-800 rounded-lg pl-10 pr-4 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
              />
           </div>
           <button className="bg-slate-900 border border-slate-800 p-2 rounded-lg hover:bg-slate-800" onClick={fetchAlerts}>
              <RefreshCw className={`h-4 w-4 ${loading ? 'animate-spin' : ''}`} />
           </button>
           <button className="bg-slate-900 border border-slate-800 p-2 rounded-lg hover:bg-slate-800">
              <Filter className="h-4 w-4" />
           </button>
        </div>
      </header>

      <div className="bg-slate-900 rounded-xl border border-slate-800 overflow-hidden">
        <table className="w-full text-left">
          <thead className="bg-slate-950 text-slate-400 text-sm uppercase font-medium">
            <tr>
              <th className="px-6 py-4">Timestamp</th>
              <th className="px-6 py-4">Detection Rule</th>
              <th className="px-6 py-4">Severity</th>
              <th className="px-6 py-4">Host</th>
              <th className="px-6 py-4">Status</th>
              <th className="px-6 py-4"></th>
            </tr>
          </thead>
          <tbody className="divide-y divide-slate-800">
            {loading ? (
              <tr>
                <td colSpan={6} className="px-6 py-12 text-center text-slate-500">Loading alerts...</td>
              </tr>
            ) : alerts.length === 0 ? (
              <tr>
                <td colSpan={6} className="px-6 py-12 text-center text-slate-500">No alerts found</td>
              </tr>
            ) : alerts.filter(a => a.rule.toLowerCase().includes(searchTerm.toLowerCase())).map((alert) => (
              <tr key={alert.id} className="hover:bg-slate-800/50 transition-colors">
                <td className="px-6 py-4 text-sm font-mono truncate max-w-[200px]">{alert.time}</td>
                <td className="px-6 py-4 font-medium">{alert.rule}</td>
                <td className="px-6 py-4">
                  <span className={`px-2 py-1 rounded text-xs font-bold uppercase ${
                    alert.severity === 'Critical' ? 'bg-red-900/50 text-red-400' :
                    alert.severity === 'High' ? 'bg-orange-900/50 text-orange-400' :
                    'bg-yellow-900/50 text-yellow-400'
                  }`}>
                    {alert.severity}
                  </span>
                </td>
                <td className="px-6 py-4 text-slate-300">{alert.host}</td>
                <td className="px-6 py-4 text-slate-400 text-sm">{alert.status}</td>
                <td className="px-6 py-4 text-right">
                  <button className="text-slate-500 hover:text-slate-200">
                    <MoreHorizontal className="h-5 w-5" />
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
