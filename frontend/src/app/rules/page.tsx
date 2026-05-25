"use client";

import { useState } from "react";
import { Shield, Plus, Search, FileCode, Play, Save } from "lucide-react";

export default function RuleManager() {
  const [rules, setRules] = useState([
    { id: "1", name: "Suspicious Logon", severity: "High", status: "Active", author: "OmniGuard" },
    { id: "2", name: "Brute Force Attempt", severity: "Medium", status: "Active", author: "OmniGuard" },
    { id: "3", name: "Process Injection", severity: "Critical", status: "Testing", author: "Analyst-01" },
  ]);

  const [selectedRule, setSelectedRule] = useState<any>(rules[0]);

  return (
    <div className="min-h-screen bg-slate-950 text-slate-50 flex">
      {/* Sidebar */}
      <div className="w-80 border-r border-slate-800 bg-slate-900/30 flex flex-col">
        <div className="p-6 border-b border-slate-800 flex justify-between items-center">
          <h2 className="font-bold flex items-center gap-2"><Shield className="h-5 w-5 text-blue-500" /> Rules</h2>
          <button className="bg-blue-600 p-1.5 rounded-lg hover:bg-blue-500"><Plus className="h-4 w-4" /></button>
        </div>
        <div className="p-4">
           <div className="relative mb-4">
              <Search className="absolute left-3 top-2.5 h-4 w-4 text-slate-500" />
              <input
                type="text"
                placeholder="Search rules..."
                className="w-full bg-slate-950 border border-slate-800 rounded-lg pl-10 pr-4 py-2 text-sm focus:outline-none focus:ring-1 focus:ring-blue-500"
              />
           </div>
           <div className="flex flex-col gap-1">
             {rules.map(rule => (
               <button
                key={rule.id}
                onClick={() => setSelectedRule(rule)}
                className={`text-left px-4 py-3 rounded-lg transition-all ${selectedRule?.id === rule.id ? 'bg-blue-600/10 text-blue-400 border border-blue-500/30' : 'hover:bg-slate-800 text-slate-400'}`}
               >
                 <div className="text-sm font-semibold">{rule.name}</div>
                 <div className="text-xs opacity-60 mt-1 uppercase tracking-tighter">{rule.severity} • {rule.status}</div>
               </button>
             ))}
           </div>
        </div>
      </div>

      {/* Editor Area */}
      <div className="flex-1 flex flex-col">
        <header className="p-6 border-b border-slate-800 flex justify-between items-center bg-slate-900/20">
          <div>
            <h1 className="text-xl font-bold">{selectedRule?.name}</h1>
            <p className="text-sm text-slate-400">Edited by {selectedRule?.author}</p>
          </div>
          <div className="flex gap-3">
             <button className="bg-slate-800 border border-slate-700 px-4 py-2 rounded-lg text-sm font-medium hover:bg-slate-700 flex items-center gap-2">
                <Play className="h-4 w-4" /> Test Rule
             </button>
             <button className="bg-blue-600 px-4 py-2 rounded-lg text-sm font-bold hover:bg-blue-500 flex items-center gap-2">
                <Save className="h-4 w-4" /> Save & Deploy
             </button>
          </div>
        </header>

        <main className="flex-1 p-8 overflow-auto">
          <div className="bg-slate-900 border border-slate-800 rounded-xl overflow-hidden shadow-2xl">
            <div className="bg-slate-950 px-4 py-2 border-b border-slate-800 flex items-center gap-2">
               <FileCode className="h-4 w-4 text-orange-400" />
               <span className="text-xs font-mono text-slate-500">sigma_rule.yml</span>
            </div>
            <textarea
              className="w-full h-[500px] bg-transparent p-6 font-mono text-sm leading-relaxed focus:outline-none resize-none text-slate-300"
              spellCheck={false}
              defaultValue={`title: ${selectedRule?.name}
id: ${selectedRule?.id}
status: experimental
description: Detects ${selectedRule?.name.toLowerCase()} activity
logsource:
    product: windows
    service: security
detection:
    selection:
        EventID: 4624
        LogonType: 3
    condition: selection
falsepositives:
    - Domain Controllers
level: high`}
            />
          </div>
        </main>
      </div>
    </div>
  );
}
