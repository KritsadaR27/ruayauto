import React, { useState, useEffect } from "react";
import { Plus, Trash2, Save } from "lucide-react";

type Pair = { keyword: string; response: string };

export default function App() {
	const [pairs, setPairs] = useState<Pair[]>([]);

	useEffect(() => {
		fetch("/api/keywords")
			.then((res) => (res.ok ? res.json() : Promise.reject()))
			.then((data) => {
				if (Array.isArray(data)) setPairs(data);
			})
			.catch(() => {});
	}, []);

	const handleChange = (
		idx: number,
		field: "keyword" | "response",
		value: string
	) => {
		setPairs((prev) =>
			prev.map((p, i) => (i === idx ? { ...p, [field]: value } : p))
		);
	};
	const handleAdd = () => setPairs([...pairs, { keyword: "", response: "" }]);
	const handleDelete = (idx: number) => setPairs(pairs.filter((_, i) => i !== idx));
	const handleSave = async () => {
		try {
			const res = await fetch("/api/keywords", {
				method: "POST",
				headers: { "Content-Type": "application/json" },
				body: JSON.stringify(pairs),
			});
			if (res.ok) {
				alert("บันทึกสำเร็จ!");
			} else {
				alert("เกิดข้อผิดพลาดในการบันทึก");
			}
		} catch (e) {
			alert("เกิดข้อผิดพลาดในการเชื่อมต่อเซิร์ฟเวอร์");
		}
	};

	return (
		<div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-sky-400 via-purple-400 to-purple-700 p-4">
			<div className="w-full max-w-2xl bg-white/90 rounded-2xl shadow-2xl p-0 overflow-hidden">
				<div className="bg-gradient-to-r from-purple-600 to-sky-400 p-6 text-center">
					<h1 className="text-3xl md:text-4xl font-bold text-white drop-shadow-lg">
						🧠 ระบบตอบกลับอัตโนมัติ (Keyword/Response)
					</h1>
				</div>
				<div className="p-6 space-y-4">
					{pairs.map((pair, idx) => (
						<div
							key={idx}
							className="bg-white/80 border-2 border-purple-300 rounded-xl shadow flex flex-col md:flex-row md:items-center gap-2 p-4 outline outline-2 outline-purple-200 transition-transform duration-200 hover:scale-[1.015] w-full"
							style={{ wordBreak: "break-word" }}
						>
							<div className="flex-1 min-w-0 flex flex-col md:flex-row md:items-center gap-2">
								<label className="text-purple-800 font-semibold md:w-28 shrink-0">
									Keyword
								</label>
								<input
									className="flex-1 border border-purple-300 rounded px-2 py-1 focus:outline-none focus:ring-2 focus:ring-purple-400 bg-white min-w-0"
									value={pair.keyword}
									onChange={(e) =>
										handleChange(idx, "keyword", e.target.value)
									}
									placeholder="เช่น hello, ราคา"
									style={{ wordBreak: "break-all" }}
								/>
							</div>
							<div className="flex-1 min-w-0 flex flex-col md:flex-row md:items-center gap-2">
								<label className="text-purple-800 font-semibold md:w-28 shrink-0">
									Response
								</label>
								<textarea
									className="flex-1 border border-purple-300 rounded px-2 py-1 focus:outline-none focus:ring-2 focus:ring-purple-400 min-h-[38px] bg-white resize-y min-w-0"
									value={pair.response}
									onChange={(e) =>
										handleChange(idx, "response", e.target.value)
									}
									placeholder="ข้อความตอบกลับ"
									style={{
										wordBreak: "break-all",
										overflowWrap: "break-word",
									}}
									rows={1}
								/>
							</div>
							<button
								className="ml-2 p-2 rounded bg-red-500 hover:bg-red-600 text-white flex items-center justify-center shrink-0"
								onClick={() => handleDelete(idx)}
								title="ลบ"
								tabIndex={-1}
							>
								<Trash2 size={20} />
							</button>
						</div>
					))}
				</div>
				<div className="flex gap-2 mt-8 justify-end p-6 pt-0 flex-col sm:flex-row space-y-2 sm:space-y-0">
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
