"use client";

import { useState, useEffect } from "react";
import {
  BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer,
  LineChart, Line, AreaChart, Area, PieChart, Pie, Cell
} from "recharts";
import { Shield, Activity, AlertTriangle, CheckCircle, Clock } from "lucide-react";

export default function SecurityDashboard() {
  const [stats, setStats] = useState({
    totalAlerts: 1254,
    criticalAlerts: 12,
    openCases: 45,
    resolvedToday: 156
  });

  const alertData = [
    { name: "00:00", count: 45 }, { name: "04:00", count: 32 },
    { name: "08:00", count: 89 }, { name: "12:00", count: 120 },
    { name: "16:00", count: 210 }, { name: "20:00", count: 150 },
  ];

  const severityData = [
    { name: "Critical", value: 12, color: "#ef4444" },
    { name: "High", value: 45, color: "#f97316" },
    { name: "Medium", value: 120, color: "#eab308" },
    { name: "Low", value: 350, color: "#3b82f6" },
  ];

  return (
    <div className="min-h-screen bg-slate-950 text-slate-50 p-8 flex flex-col gap-8">
      <header>
        <h1 className="text-3xl font-bold flex items-center gap-3">
          <Shield className="text-blue-500 h-8 w-8" /> OmniGuard Command Center
        </h1>
        <p className="text-slate-400 mt-2">Real-time visibility into your security posture</p>
      </header>

      <div className="grid grid-cols-1 md:grid-cols-4 gap-6">
        <StatCard title="Total Events" value={stats.totalAlerts.toLocaleString()} icon={<Activity className="text-blue-400" />} trend="+12%" />
        <StatCard title="Critical Detections" value={stats.criticalAlerts} icon={<AlertTriangle className="text-red-400" />} trend="+2" />
        <StatCard title="Active Investigations" value={stats.openCases} icon={<Clock className="text-orange-400" />} trend="-4" />
        <StatCard title="Resolved (24h)" value={stats.resolvedToday} icon={<CheckCircle className="text-green-400" />} trend="+24" />
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
        <div className="bg-slate-900/50 border border-slate-800 p-6 rounded-2xl">
          <h2 className="text-lg font-semibold mb-6">Ingestion Volume (EPS)</h2>
          <div className="h-[300px]">
            <ResponsiveContainer width="100%" height="100%">
              <AreaChart data={alertData}>
                <defs>
                  <linearGradient id="colorCount" x1="0" y1="0" x2="0" y2="1">
                    <stop offset="5%" stopColor="#3b82f6" stopOpacity={0.3}/>
                    <stop offset="95%" stopColor="#3b82f6" stopOpacity={0}/>
                  </linearGradient>
                </defs>
                <CartesianGrid strokeDasharray="3 3" stroke="#1e293b" />
                <XAxis dataKey="name" stroke="#64748b" />
                <YAxis stroke="#64748b" />
                <Tooltip
                  contentStyle={{ backgroundColor: "#0f172a", border: "1px solid #1e293b" }}
                  itemStyle={{ color: "#3b82f6" }}
                />
                <Area type="monotone" dataKey="count" stroke="#3b82f6" fillOpacity={1} fill="url(#colorCount)" strokeWidth={2} />
              </AreaChart>
            </ResponsiveContainer>
          </div>
        </div>

        <div className="bg-slate-900/50 border border-slate-800 p-6 rounded-2xl">
          <h2 className="text-lg font-semibold mb-6">Alert Severity Distribution</h2>
          <div className="h-[300px] flex items-center">
            <ResponsiveContainer width="100%" height="100%">
              <PieChart>
                <Pie
                  data={severityData}
                  cx="50%"
                  cy="50%"
                  innerRadius={60}
                  outerRadius={100}
                  paddingAngle={5}
                  dataKey="value"
                >
                  {severityData.map((entry, index) => (
                    <Cell key={`cell-${index}`} fill={entry.color} />
                  ))}
                </Pie>
                <Tooltip
                  contentStyle={{ backgroundColor: "#0f172a", border: "1px solid #1e293b" }}
                />
              </PieChart>
            </ResponsiveContainer>
            <div className="flex flex-col gap-3 pr-8">
               {severityData.map(s => (
                 <div key={s.name} className="flex items-center gap-2">
                    <div className="w-3 h-3 rounded-full" style={{backgroundColor: s.color}} />
                    <span className="text-sm text-slate-400">{s.name}: <span className="text-slate-100 font-mono">{s.value}</span></span>
                 </div>
               ))}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

function StatCard({ title, value, icon, trend }: { title: string, value: string | number, icon: React.ReactNode, trend: string }) {
  const isPositive = trend.startsWith("+");
  return (
    <div className="bg-slate-900 border border-slate-800 p-6 rounded-2xl hover:border-slate-700 transition-all">
      <div className="flex justify-between items-start mb-4">
        <div className="p-2 bg-slate-950 rounded-lg">{icon}</div>
        <span className={`text-xs font-bold px-2 py-1 rounded ${isPositive ? 'bg-green-900/30 text-green-400' : 'bg-red-900/30 text-red-400'}`}>
          {trend}
        </span>
      </div>
      <h3 className="text-slate-400 text-sm font-medium">{title}</h3>
      <p className="text-2xl font-bold mt-1">{value}</p>
    </div>
  );
}
