"use client";

import { AlertTrends } from "@/components/dashboard/AlertTrends";
import { Shield, AlertTriangle, Activity, Users } from "lucide-react";

export default function Dashboard() {
  return (
    <main className="min-h-screen bg-slate-950 text-slate-50 p-8">
      <header className="mb-8">
        <h1 className="text-3xl font-bold">SOC Overview</h1>
        <p className="text-slate-400">OmniGuard Security Operations Center Dashboard</p>
      </header>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
        <StatCard title="Total Alerts" value="1,284" icon={<AlertTriangle className="text-amber-500" />} />
        <StatCard title="Active Cases" value="12" icon={<Shield className="text-blue-500" />} />
        <StatCard title="EPS (Events/sec)" value="4.2k" icon={<Activity className="text-emerald-500" />} />
        <StatCard title="Online Agents" value="452" icon={<Users className="text-slate-400" />} />
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
        <AlertTrends />
        <div className="bg-slate-900 p-6 rounded-xl border border-slate-800">
           <h3 className="text-slate-200 font-semibold mb-4">Top Detection Rules</h3>
           <div className="space-y-4">
              <RuleItem name="Suspicious PowerShell Execution" count={142} />
              <RuleItem name="Brute Force Attempt" count={85} />
              <RuleItem name="Anomalous Domain Query" count={34} />
           </div>
        </div>
      </div>
    </main>
  );
}

function StatCard({ title, value, icon }: { title: string, value: string, icon: React.ReactNode }) {
  return (
    <div className="bg-slate-900 p-6 rounded-xl border border-slate-800 flex items-center justify-between">
      <div>
        <p className="text-slate-400 text-sm font-medium">{title}</p>
        <p className="text-2xl font-bold mt-1">{value}</p>
      </div>
      <div className="p-3 bg-slate-950 rounded-lg border border-slate-800">
        {icon}
      </div>
    </div>
  );
}

function RuleItem({ name, count }: { name: string, count: number }) {
  return (
    <div className="flex items-center justify-between">
      <span className="text-slate-300">{name}</span>
      <span className="text-slate-500 font-mono">{count}</span>
    </div>
  );
}
