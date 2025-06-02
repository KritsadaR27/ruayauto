import React, { useState, useEffect, useRef } from "react";
import { Plus, Trash2, Save } from "lucide-react";

type Pair = { keywords: string[]; responses: string[] };

export default function App() {
	const [pairs, setPairs] = useState<Pair[]>([]);
	const [defaultResponses, setDefaultResponses] = useState<string[]>([""]);
	const [enableDefault, setEnableDefault] = useState(true);
	const [noTag, setNoTag] = useState(true);
	const [noSticker, setNoSticker] = useState(true);
	const keywordRefs = useRef<{ [key: string]: HTMLInputElement | null }>({});

	useEffect(() => {
		fetch("/api/keywords")
			.then((res) => (res.ok ? res.json() : Promise.reject()))
			.then((data) => {
				if (data && typeof data === "object") {
					if (Array.isArray(data.pairs)) setPairs(data.pairs);
					else if (Array.isArray(data)) setPairs(data); // fallback
					if (Array.isArray(data.defaultResponses)) setDefaultResponses(data.defaultResponses);
					if (typeof data.enableDefault === "boolean") setEnableDefault(data.enableDefault);
					if (typeof data.noTag === "boolean") setNoTag(data.noTag);
					if (typeof data.noSticker === "boolean") setNoSticker(data.noSticker);
				}
			})
			.catch(() => {});
	}, []);

	// เพิ่ม/ลบ/แก้ไข keyword/response ในแต่ละ pair
	const handleKeywordChange = (pairIdx: number, kwIdx: number, value: string) => {
		setPairs((prev) =>
			prev.map((p, i) =>
				i === pairIdx
					? { ...p, keywords: p.keywords.map((k, j) => (j === kwIdx ? value : k)) }
					: p
			)
		);
	};
	const handleAddKeyword = (pairIdx: number) => {
		setPairs((prev) =>
			prev.map((p, i) =>
				i === pairIdx ? { ...p, keywords: [...p.keywords, ""] } : p
			)
		);
		// Focus on the new keyword input after state update
		setTimeout(() => {
			const newKeywordIndex = pairs[pairIdx].keywords.length;
			const refKey = `${pairIdx}-${newKeywordIndex}`;
			keywordRefs.current[refKey]?.focus();
		}, 0);
	};
	const handleDeleteKeyword = (pairIdx: number, kwIdx: number) => {
		setPairs((prev) =>
			prev.map((p, i) =>
				i === pairIdx
					? { ...p, keywords: p.keywords.filter((_, j) => j !== kwIdx) }
					: p
			)
		);
	};
	const handleResponseChange = (pairIdx: number, resIdx: number, value: string) => {
		setPairs((prev) =>
			prev.map((p, i) =>
				i === pairIdx
					? { ...p, responses: p.responses.map((r, j) => (j === resIdx ? value : r)) }
					: p
			)
		);
	};
	const handleAddResponse = (pairIdx: number) => {
		setPairs((prev) =>
			prev.map((p, i) =>
				i === pairIdx ? { ...p, responses: [...p.responses, ""] } : p
			)
		);
	};
	const handleDeleteResponse = (pairIdx: number, resIdx: number) => {
		setPairs((prev) =>
			prev.map((p, i) =>
				i === pairIdx
					? { ...p, responses: p.responses.filter((_, j) => j !== resIdx) }
					: p
			)
		);
	};
	const handleAddPair = () => setPairs([...pairs, { keywords: [""], responses: [""] }]);
	const handleDeletePair = (idx: number) => setPairs(pairs.filter((_, i) => i !== idx));
	const handleDefaultResponseChange = (idx: number, value: string) => {
		setDefaultResponses((prev) => prev.map((r, i) => (i === idx ? value : r)));
	};
	const handleAddDefaultResponse = () => setDefaultResponses([...defaultResponses, ""]);
	const handleDeleteDefaultResponse = (idx: number) => setDefaultResponses(defaultResponses.filter((_, i) => i !== idx));

	const handleSave = async () => {
		// ส่งข้อมูลแบบใหม่ไป backend
		try {
			const res = await fetch("/api/keywords", {
				method: "POST",
				headers: { "Content-Type": "application/json" },
				body: JSON.stringify({
					pairs,
					defaultResponses,
					enableDefault,
					noTag,
					noSticker,
				}),
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
		<div className="min-h-screen bg-gradient-to-br from-blue-50 via-indigo-50 to-purple-50 p-4">
			<div className="max-w-4xl mx-auto">
				{/* Header Section */}
				<div className="bg-white rounded-2xl shadow-lg border border-gray-100 mb-6 overflow-hidden">
					<div className="bg-gradient-to-r from-indigo-600 via-purple-600 to-blue-600 px-8 py-6">
						<div className="flex items-center justify-between">
							<div>
								<h1 className="text-2xl font-bold text-white mb-2">
									🤖 ระบบตอบกลับอัตโนมัติ Facebook
								</h1>
								<p className="text-indigo-100 text-sm">
									จัดการคีย์เวิร์ดและข้อความตอบกลับ
								</p>
							</div>
							<div className="hidden md:flex items-center space-x-4">
								<div className="bg-white/10 rounded-lg px-4 py-2">
									<span className="text-white text-sm font-medium">
										{pairs.length} รายการ
									</span>
								</div>
							</div>
						</div>
					</div>
				</div>

				{/* Settings Section */}
				<div className="bg-white rounded-2xl shadow-lg border border-gray-100 mb-6 overflow-hidden">
					<div className="px-8 py-6 border-b border-gray-100">
						<h2 className="text-lg font-semibold text-gray-900 mb-4">
							⚙️ การตั้งค่าระบบ
						</h2>
						
						<div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
							{/* Default Response Settings */}
							<div className="space-y-4">
								<div className="flex items-start space-x-3">
									<input 
										type="checkbox" 
										id="enableDefault"
										className="mt-1 h-4 w-4 text-indigo-600 focus:ring-indigo-500 border-gray-300 rounded"
										checked={enableDefault} 
										onChange={e => setEnableDefault(e.target.checked)} 
									/>
									<div className="flex-1">
										<label htmlFor="enableDefault" className="text-sm font-medium text-gray-900 cursor-pointer">
											คำตอบกรณีไม่ตรงคีย์เวิร์ด
										</label>
										<p className="text-xs text-gray-500 mt-1">
											ระบบจะสุ่มเลือกคำตอบจากรายการด้านล่าง
										</p>
									</div>
								</div>
								
								{enableDefault && (
									<div className="ml-7 space-y-3">
										{defaultResponses.map((res, idx) => (
											<div key={idx} className="flex items-start space-x-2">
												<textarea
													className="flex-1 px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 focus:shadow-lg resize-none text-sm transition-all duration-200"
													value={res}
													onChange={e => handleDefaultResponseChange(idx, e.target.value)}
													placeholder="ข้อความตอบกลับเมื่อไม่ตรงคีย์เวิร์ด"
													rows={2}
												/>
												{defaultResponses.length > 1 && (
													<button 
														className="mt-1 p-1.5 text-red-500 hover:text-red-700 hover:bg-red-50 rounded-lg transition-colors" 
														onClick={() => handleDeleteDefaultResponse(idx)} 
														type="button"
													>
														<Trash2 size={16} />
													</button>
												)}
											</div>
										))}
										<button 
											className="inline-flex items-center px-3 py-1.5 text-xs font-medium text-indigo-700 bg-indigo-100 hover:bg-indigo-200 rounded-lg transition-colors" 
											onClick={handleAddDefaultResponse} 
											type="button"
										>
											<Plus size={14} className="mr-1" />
											เพิ่มข้อความ
										</button>
									</div>
								)}
							</div>

							{/* Comment Filter Settings */}
							<div className="space-y-4">
								<h3 className="text-sm font-medium text-gray-900">
									🛡️ ตัวกรองคอมเมนต์
								</h3>
								<div className="space-y-3">
									<label className="flex items-center space-x-3 cursor-pointer group">
										<input 
											type="checkbox" 
											className="h-4 w-4 text-indigo-600 focus:ring-indigo-500 border-gray-300 rounded"
											checked={noTag} 
											onChange={e => setNoTag(e.target.checked)} 
										/>
										<div>
											<span className="text-sm font-medium text-gray-900 group-hover:text-indigo-600 transition-colors">
												ไม่ตอบคอมเมนต์ที่มีการแท็ก
											</span>
											<p className="text-xs text-gray-500">
												ข้ามคอมเมนต์ที่มี @ mention
											</p>
										</div>
									</label>
									<label className="flex items-center space-x-3 cursor-pointer group">
										<input 
											type="checkbox" 
											className="h-4 w-4 text-indigo-600 focus:ring-indigo-500 border-gray-300 rounded"
											checked={noSticker} 
											onChange={e => setNoSticker(e.target.checked)} 
										/>
										<div>
											<span className="text-sm font-medium text-gray-900 group-hover:text-indigo-600 transition-colors">
												ไม่ตอบคอมเมนต์ที่มีสติกเกอร์
											</span>
											<p className="text-xs text-gray-500">
												ข้ามคอมเมนต์ที่มี emoji/sticker
											</p>
										</div>
									</label>
								</div>
							</div>
						</div>
					</div>
				</div>

				{/* Keyword-Response Pairs */}
				<div className="space-y-4">
					{pairs.map((pair, idx) => (
						<div
							key={idx}
							className="bg-white rounded-2xl shadow-lg border border-gray-100 overflow-hidden hover:shadow-xl transition-all duration-300"
						>
							<div className="px-8 py-6">
								<div className="flex items-center justify-between mb-6">
									<h3 className="text-lg font-semibold text-gray-900">
										📝 รายการที่ {idx + 1}
									</h3>
									<button
										className="p-2 text-red-500 hover:text-red-700 hover:bg-red-50 rounded-lg transition-colors"
										onClick={() => handleDeletePair(idx)}
										title="ลบรายการนี้"
										type="button"
									>
										<Trash2 size={18} />
									</button>
								</div>

								<div className="space-y-6">
									{/* Keywords Section */}
									<div className="space-y-4">
										<div className="flex items-center justify-between">
											<h4 className="text-sm font-semibold text-gray-900 flex items-center">
												🔍 คีย์เวิร์ด
											</h4>
											<button 
												className="inline-flex items-center px-3 py-1.5 text-xs font-medium text-indigo-700 bg-indigo-100 hover:bg-indigo-200 rounded-lg transition-colors" 
												onClick={() => handleAddKeyword(idx)} 
												type="button"
											>
												<Plus size={14} className="mr-1" />
												เพิ่ม
											</button>
										</div>
										<div className="flex flex-wrap gap-2">
											{pair.keywords.map((kw, kidx) => (
												<div key={kidx} className="inline-flex items-center bg-indigo-50 border border-indigo-200 rounded-full px-3 py-2 group hover:bg-indigo-100 transition-colors focus-within:bg-indigo-100 focus-within:border-indigo-400 focus-within:shadow-md">
													<input
														ref={(el) => {
															keywordRefs.current[`${idx}-${kidx}`] = el;
														}}
														className="bg-transparent border-none outline-none text-sm font-medium text-indigo-800 min-w-[60px] max-w-[120px] text-center focus:ring-0 focus:outline-none"
														value={kw}
														onChange={e => handleKeywordChange(idx, kidx, e.target.value)}
														placeholder="คีย์เวิร์ด"
													/>
													{pair.keywords.length > 1 && (
														<button 
															className="ml-2 text-indigo-400 hover:text-red-500 transition-colors" 
															onClick={() => handleDeleteKeyword(idx, kidx)} 
															type="button"
														>
															×
														</button>
													)}
												</div>
											))}
										</div>
									</div>

									{/* Responses Section */}
									<div className="space-y-4">
										<div className="flex items-center justify-between">
											<h4 className="text-sm font-semibold text-gray-900 flex items-center">
												💬 คำตอบ
											</h4>
											<button 
												className="inline-flex items-center px-3 py-1.5 text-xs font-medium text-purple-700 bg-purple-100 hover:bg-purple-200 rounded-lg transition-colors" 
												onClick={() => handleAddResponse(idx)} 
												type="button"
											>
												<Plus size={14} className="mr-1" />
												เพิ่ม
											</button>
										</div>
										<div className="space-y-3">
											{pair.responses.map((res, ridx) => (
												<div key={ridx} className="flex items-start space-x-2">
													<textarea
														className="flex-1 px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-purple-500 focus:shadow-lg resize-none text-sm transition-all duration-200"
														value={res}
														onChange={e => handleResponseChange(idx, ridx, e.target.value)}
														placeholder="ข้อความตอบกลับ"
														rows={3}
													/>
													{pair.responses.length > 1 && (
														<button 
															className="mt-1 p-1.5 text-red-500 hover:text-red-700 hover:bg-red-50 rounded-lg transition-colors" 
															onClick={() => handleDeleteResponse(idx, ridx)} 
															type="button"
														>
															<Trash2 size={16} />
														</button>
													)}
												</div>
											))}
										</div>
									</div>
								</div>
							</div>
						</div>
					))}
				</div>

				{/* Action Buttons */}
				<div className="flex flex-col sm:flex-row gap-4 items-center justify-between mt-8 bg-white rounded-2xl shadow-lg border border-gray-100 px-8 py-6">
					<div className="text-sm text-gray-600">
						💡 เคล็ดลับ: สามารถเพิ่มคีย์เวิร์ดและคำตอบหลายๆ อันในแต่ละรายการ
					</div>
					<div className="flex gap-3">
						<button
							className="inline-flex items-center px-6 py-3 bg-gradient-to-r from-indigo-500 to-purple-600 hover:from-indigo-600 hover:to-purple-700 text-white font-semibold rounded-xl shadow-lg hover:shadow-xl transition-all duration-300"
							onClick={handleAddPair}
							type="button"
						>
							<Plus size={20} className="mr-2" />
							เพิ่มรายการใหม่
						</button>
						<button
							className="inline-flex items-center px-6 py-3 bg-gradient-to-r from-green-500 to-emerald-600 hover:from-green-600 hover:to-emerald-700 text-white font-semibold rounded-xl shadow-lg hover:shadow-xl transition-all duration-300"
							onClick={handleSave}
							type="button"
						>
							<Save size={20} className="mr-2" />
							บันทึกการเปลี่ยนแปลง
						</button>
					</div>
				</div>
			</div>
		</div>
	);
}
