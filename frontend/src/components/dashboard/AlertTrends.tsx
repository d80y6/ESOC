"use client";

import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from "recharts";

const data = [
  { name: "Mon", alerts: 40 },
  { name: "Tue", alerts: 30 },
  { name: "Wed", alerts: 20 },
  { name: "Thu", alerts: 27 },
  { name: "Fri", alerts: 18 },
  { name: "Sat", alerts: 23 },
  { name: "Sun", alerts: 34 },
];

export const AlertTrends = () => {
  return (
    <div className="h-[300px] w-full bg-slate-900 p-4 rounded-xl border border-slate-800">
      <h3 className="text-slate-200 font-semibold mb-4">Alert Trends (Last 7 Days)</h3>
      <ResponsiveContainer width="100%" height="80%">
        <BarChart data={data}>
          <CartesianGrid strokeDasharray="3 3" stroke="#334155" />
          <XAxis dataKey="name" stroke="#94a3b8" />
          <YAxis stroke="#94a3b8" />
          <Tooltip
            contentStyle={{ backgroundColor: "#0f172a", border: "1px solid #1e293b" }}
            itemStyle={{ color: "#f1f5f9" }}
          />
          <Bar dataKey="alerts" fill="#3b82f6" radius={[4, 4, 0, 0]} />
        </BarChart>
      </ResponsiveContainer>
    </div>
  );
};
