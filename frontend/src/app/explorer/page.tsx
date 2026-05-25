"use client";

import { useState } from "react";
import { Search, Filter, Play, Download, Trash2, ShieldAlert } from "lucide-react";

export default function LogExplorer() {
  const [query, setQuery] = useState("*");
  const [results, setResults] = useState<any[]>([]);
  const [loading, setLoading] = useState(false);

  const handleSearch = async () => {
    setLoading(true);
    try {
      const response = await fetch("/api/search", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ query, size: 50 }),
      });
      const data = await response.json();
      setResults(data.hits?.hits || []);
    } catch (err) {
      console.error("Search failed", err);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-slate-950 text-slate-50 p-8 flex flex-col gap-6">
      <header className="flex justify-between items-center">
        <div>
          <h1 className="text-3xl font-bold flex items-center gap-2">
            <ShieldAlert className="text-blue-500" /> Log Explorer
          </h1>
          <p className="text-slate-400">Hunt for threats across all normalized logs</p>
        </div>
        <div className="flex gap-2">
           <button className="bg-slate-900 border border-slate-800 px-4 py-2 rounded-lg text-sm flex items-center gap-2 hover:bg-slate-800">
              <Download className="h-4 w-4" /> Export CSV
           </button>
           <button className="bg-blue-600 px-4 py-2 rounded-lg text-sm font-semibold flex items-center gap-2 hover:bg-blue-500" onClick={handleSearch}>
              <Play className="h-4 w-4" /> Run Query
           </button>
        </div>
      </header>

      <div className="flex flex-col gap-4">
        <div className="bg-slate-900 border border-slate-800 p-4 rounded-xl flex gap-4">
          <div className="flex-1 relative">
            <Search className="absolute left-3 top-3 h-5 w-5 text-slate-500" />
            <input
              type="text"
              className="w-full bg-slate-950 border border-slate-800 rounded-lg pl-12 pr-4 py-3 text-lg font-mono focus:outline-none focus:ring-2 focus:ring-blue-500"
              placeholder="event.action: logon AND user.name: admin..."
              value={query}
              onChange={(e) => setQuery(e.target.value)}
            />
          </div>
          <button className="bg-slate-800 border border-slate-700 p-3 rounded-lg hover:bg-slate-700">
            <Filter className="h-5 w-5" />
          </button>
        </div>

        <div className="bg-slate-900 rounded-xl border border-slate-800 min-h-[400px] overflow-hidden">
          {loading ? (
             <div className="flex items-center justify-center h-64">
                <div className="animate-spin rounded-full h-8 w-8 border-t-2 border-b-2 border-blue-500"></div>
             </div>
          ) : results.length > 0 ? (
            <table className="w-full text-left">
              <thead className="bg-slate-950 text-slate-400 text-xs uppercase font-bold tracking-wider">
                <tr>
                  <th className="px-6 py-4">Timestamp</th>
                  <th className="px-6 py-4">Source</th>
                  <th className="px-6 py-4">Action</th>
                  <th className="px-6 py-4">User / Host</th>
                  <th className="px-6 py-4">Message</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-slate-800 font-mono text-sm">
                {results.map((hit, i) => (
                  <tr key={i} className="hover:bg-slate-800/50 transition-colors">
                    <td className="px-6 py-4 text-blue-400 whitespace-nowrap">{hit._source["@timestamp"]}</td>
                    <td className="px-6 py-4">{hit._source["event.provider"] || hit._source["log.file.path"]}</td>
                    <td className="px-6 py-4">
                       <span className="bg-slate-800 px-2 py-0.5 rounded text-slate-300 border border-slate-700">
                          {hit._source["event.action"] || "N/A"}
                       </span>
                    </td>
                    <td className="px-6 py-4">
                       <div className="flex flex-col">
                          <span className="text-slate-200">{hit._source["user.name"] || "system"}</span>
                          <span className="text-xs text-slate-500">{hit._source["host.name"]}</span>
                       </div>
                    </td>
                    <td className="px-6 py-4 max-w-md truncate text-slate-400">{hit._source["message"]}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          ) : (
            <div className="flex flex-col items-center justify-center h-64 text-slate-500">
               <Trash2 className="h-12 w-12 mb-4 opacity-20" />
               <p>No logs found for this query</p>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
