"use client";

import { Shield, Clock, FileText, Share2, Plus } from "lucide-react";

export default function CaseWorkbench({ params }: { params: { id: string } }) {
  return (
    <div className="min-h-screen bg-slate-950 text-slate-50 flex flex-col">
      <header className="border-b border-slate-800 p-6 flex justify-between items-center bg-slate-900/50">
        <div className="flex items-center gap-4">
          <div className="bg-blue-900/30 p-2 rounded-lg">
            <Shield className="text-blue-500 h-6 w-6" />
          </div>
          <div>
            <h1 className="text-xl font-bold">Case #{params.id}: APT29 Lateral Movement</h1>
            <p className="text-slate-400 text-sm">Assigned to: <span className="text-blue-400">jules_engineer</span></p>
          </div>
        </div>
        <div className="flex gap-3">
          <button className="flex items-center gap-2 bg-slate-800 px-4 py-2 rounded-lg text-sm font-medium hover:bg-slate-700">
            <Share2 className="h-4 w-4" /> Share
          </button>
          <button className="flex items-center gap-2 bg-blue-600 px-4 py-2 rounded-lg text-sm font-medium hover:bg-blue-500 transition-colors">
            Resolve Case
          </button>
        </div>
      </header>

      <div className="flex-1 grid grid-cols-12 overflow-hidden">
        {/* Timeline */}
        <div className="col-span-8 p-8 border-r border-slate-800 overflow-y-auto">
          <h2 className="text-lg font-bold mb-6 flex items-center gap-2">
            <Clock className="h-5 w-5 text-slate-400" /> Incident Timeline
          </h2>
          <div className="space-y-8 relative before:absolute before:left-2 before:top-2 before:bottom-2 before:w-px before:bg-slate-800">
            <TimelineItem
              time="14:20"
              title="Malicious Document Opened"
              desc="Winword.exe spawned suspicious child process cmd.exe on WRK-99"
              type="Alert"
            />
            <TimelineItem
              time="14:25"
              title="External Connection Established"
              desc="Connection to 194.165.16.2 over port 443 detected"
              type="Network"
            />
            <TimelineItem
              time="14:30"
              title="Lateral Movement Detected"
              desc="SMB connection from WRK-99 to DC-01 using Domain Admin credentials"
              type="Identity"
            />
          </div>
        </div>

        {/* Evidence & Sidebar */}
        <div className="col-span-4 p-8 bg-slate-900/30 overflow-y-auto">
           <div className="mb-10">
              <div className="flex justify-between items-center mb-4">
                <h2 className="text-lg font-bold flex items-center gap-2">
                  <FileText className="h-5 w-5 text-slate-400" /> Evidence
                </h2>
                <button className="text-blue-500 hover:text-blue-400">
                  <Plus className="h-4 w-4" />
                </button>
              </div>
              <div className="space-y-3">
                <EvidenceCard type="IP Address" value="194.165.16.2" tags={["C2", "Known Malicious"]} />
                <EvidenceCard type="File Hash" value="7e52b610c1...d9f" tags={["Cobalt Strike"]} />
                <EvidenceCard type="User Account" value="svc_backup" tags={["Compromised"]} />
              </div>
           </div>

           <div>
              <h2 className="text-lg font-bold mb-4">Case Summary</h2>
              <div className="bg-slate-900 p-4 rounded-lg border border-slate-800 text-sm text-slate-300 leading-relaxed">
                Automated analysis indicates a potential multi-stage attack starting with a phishing document.
                Lateral movement to the Domain Controller was observed 10 minutes after initial access.
              </div>
           </div>
        </div>
      </div>
    </div>
  );
}

function TimelineItem({ time, title, desc, type }: { time: string, title: string, desc: string, type: string }) {
  return (
    <div className="pl-10 relative">
      <div className="absolute left-0 top-1.5 w-4 h-4 bg-slate-950 border-2 border-blue-500 rounded-full z-10" />
      <div className="text-xs font-mono text-slate-500 mb-1">{time}</div>
      <div className="bg-slate-900 border border-slate-800 p-4 rounded-xl">
        <div className="flex justify-between items-center mb-2">
          <span className="font-bold">{title}</span>
          <span className="text-[10px] px-1.5 py-0.5 bg-slate-800 rounded uppercase tracking-wider font-bold text-slate-400">{type}</span>
        </div>
        <p className="text-sm text-slate-400">{desc}</p>
      </div>
    </div>
  );
}

function EvidenceCard({ type, value, tags }: { type: string, value: string, tags: string[] }) {
  return (
    <div className="bg-slate-900 border border-slate-800 p-4 rounded-xl hover:border-slate-700 transition-colors cursor-pointer">
      <div className="text-[10px] text-slate-500 uppercase font-bold mb-1">{type}</div>
      <div className="font-mono text-sm mb-2 break-all">{value}</div>
      <div className="flex gap-2">
        {tags.map(tag => (
          <span key={tag} className="text-[10px] px-1.5 py-0.5 bg-blue-900/20 text-blue-400 rounded border border-blue-900/50">{tag}</span>
        ))}
      </div>
    </div>
  );
}
