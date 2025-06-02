import React, { useState } from "react";
import { Plus, Trash2, Save } from "lucide-react";

const initialPairs = [
  { keyword: "hello", response: "Hi there!" },
  { keyword: "help", response: "How can I assist you?" },
  { keyword: "ราคา", response: "ราคาสินค้าเริ่มต้นที่ 99 บาทค่ะ" },
  { keyword: "ส่งฟรีไหม", response: "เรามีบริการส่งฟรีเมื่อยอดครบ 500 บาทค่ะ" },
];

export default function App() {
  const [pairs, setPairs] = useState(initialPairs);

  const handleChange = (idx: number, field: "keyword" | "response", value: string) => {
    setPairs((prev) =>
      prev.map((p, i) => (i === idx ? { ...p, [field]: value } : p))
    );
  };

  const handleAdd = () => setPairs([...pairs, { keyword: "", response: "" }]);
  const handleDelete = (idx: number) => setPairs(pairs.filter((_, i) => i !== idx));
  const handleSave = () => {
    // TODO: Save logic (API call or export)
    alert("บันทึกสำเร็จ!");
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-sky-400 via-purple-500 to-purple-700 p-4">
      <div className="w-full max-w-2xl">
        <div className="mb-6 text-center">
          <h1 className="text-3xl md:text-4xl font-bold text-white drop-shadow-lg">
            🧠 ระบบตอบกลับอัตโนมัติ
          </h1>
        </div>
        <div className="space-y-4">
          {pairs.map((pair, idx) => (
            <div
              key={idx}
              className="bg-white/80 border-2 border-purple-300 rounded-xl shadow flex flex-col md:flex-row md:items-center gap-2 p-4 outline outline-2 outline-purple-200"
            >
              <div className="flex-1 flex flex-col md:flex-row md:items-center gap-2">
                <label className="text-purple-800 font-semibold md:w-28">Keyword</label>
                <input
                  className="flex-1 border border-purple-300 rounded px-2 py-1 focus:outline-none focus:ring-2 focus:ring-purple-400"
                  value={pair.keyword}
                  onChange={e => handleChange(idx, "keyword", e.target.value)}
                  placeholder="เช่น hello, ราคา"
                />
              </div>
              <div className="flex-1 flex flex-col md:flex-row md:items-center gap-2">
                <label className="text-purple-800 font-semibold md:w-28">Response</label>
                <textarea
                  className="flex-1 border border-purple-300 rounded px-2 py-1 focus:outline-none focus:ring-2 focus:ring-purple-400 min-h-[38px]"
                  value={pair.response}
                  onChange={e => handleChange(idx, "response", e.target.value)}
                  placeholder="ข้อความตอบกลับ"
                />
              </div>
              <button
                className="ml-2 p-2 rounded bg-red-500 hover:bg-red-600 text-white flex items-center justify-center"
                onClick={() => handleDelete(idx)}
                title="ลบ"
              >
                <Trash2 size={20} />
              </button>
            </div>
          ))}
        </div>
        <div className="flex gap-2 mt-8 justify-end">
          <button
            className="flex items-center gap-1 px-4 py-2 rounded bg-yellow-400 hover:bg-yellow-500 text-purple-900 font-bold shadow"
            onClick={handleAdd}
          >
            <Plus size={20} /> เพิ่มใหม่
          </button>
          <button
            className="flex items-center gap-1 px-4 py-2 rounded bg-green-500 hover:bg-green-600 text-white font-bold shadow"
            onClick={handleSave}
          >
            <Save size={20} /> บันทึก
          </button>
        </div>
      </div>
    </div>
  );
}
