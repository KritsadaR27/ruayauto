import { createSignal, For, onMount } from "solid-js";

interface Pair {
  keyword: string;
  response: string;
}

const initialPairs: Pair[] = [
  { keyword: "hello", response: "Hi there!" },
  { keyword: "help", response: "How can I assist you?" },
  { keyword: "ราคา", response: "ราคาสินค้าเริ่มต้นที่ 99 บาทค่ะ" },
  { keyword: "ส่งฟรีไหม", response: "เรามีบริการส่งฟรีเมื่อยอดครบ 500 บาทค่ะ" },
];

export default function App() {
  const [pairs, setPairs] = createSignal<Pair[]>(initialPairs);

  onMount(async () => {
    try {
      const res = await fetch("/api/keywords");
      if (res.ok) {
        const data = await res.json();
        if (Array.isArray(data)) setPairs(data);
      }
    } catch {}
  });

  const handleChange = (idx: number, field: keyof Pair, value: string) => {
    setPairs((prev) => prev.map((p, i) => (i === idx ? { ...p, [field]: value } : p)));
  };

  const handleAdd = () => setPairs([...pairs(), { keyword: "", response: "" }]);
  const handleDelete = (idx: number) => setPairs(pairs().filter((_, i) => i !== idx));
  const handleSave = async () => {
    try {
      const res = await fetch("/api/keywords", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(pairs()),
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
    <div style={{ "min-height": "100vh", display: "flex", "align-items": "center", "justify-content": "center", background: "linear-gradient(135deg, #38bdf8, #a78bfa 60%, #7c3aed)" }}>
      <div style={{ width: "100%", "max-width": "600px" }}>
        <div style={{ margin: "1.5rem 0", "text-align": "center" }}>
          <h1 style={{ "font-size": "2rem", color: "#fff", "font-weight": "bold", "text-shadow": "0 2px 8px #0002" }}>
            🧠 ระบบตอบกลับอัตโนมัติ (SolidJS)
          </h1>
        </div>
        <div style={{ "margin-bottom": "2rem" }}>
          <For each={pairs()}>{(pair, idx) => (
            <div style={{ background: "#fff8", border: "2px solid #a78bfa", "border-radius": "1rem", margin: "0.5rem 0", padding: "1rem", display: "flex", gap: "0.5rem", "align-items": "center" }}>
              <div style={{ flex: 1 }}>
                <label style={{ color: "#6d28d9", "font-weight": "bold" }}>Keyword</label>
                <input
                  style={{ width: "100%", border: "1px solid #a78bfa", "border-radius": "0.5rem", padding: "0.25rem 0.5rem" }}
                  value={pair.keyword}
                  onInput={e => handleChange(idx(), "keyword", e.currentTarget.value)}
                  placeholder="เช่น hello, ราคา"
                />
              </div>
              <div style={{ flex: 1 }}>
                <label style={{ color: "#6d28d9", "font-weight": "bold" }}>Response</label>
                <textarea
                  style={{ width: "100%", border: "1px solid #a78bfa", "border-radius": "0.5rem", padding: "0.25rem 0.5rem", "min-height": "38px" }}
                  value={pair.response}
                  onInput={e => handleChange(idx(), "response", e.currentTarget.value)}
                  placeholder="ข้อความตอบกลับ"
                />
              </div>
              <button
                style={{ marginLeft: "0.5rem", padding: "0.5rem", background: "#ef4444", color: "#fff", border: "none", "border-radius": "0.5rem", cursor: "pointer" }}
                onClick={() => handleDelete(idx())}
                title="ลบ"
              >
                🗑️
              </button>
            </div>
          )}</For>
        </div>
        <div style={{ display: "flex", gap: "0.5rem", "justify-content": "flex-end" }}>
          <button
            style={{ padding: "0.5rem 1rem", background: "#fde047", color: "#7c3aed", "font-weight": "bold", border: "none", "border-radius": "0.5rem", cursor: "pointer" }}
            onClick={handleAdd}
          >
            ➕ เพิ่มใหม่
          </button>
          <button
            style={{ padding: "0.5rem 1rem", background: "#22c55e", color: "#fff", "font-weight": "bold", border: "none", "border-radius": "0.5rem", cursor: "pointer" }}
            onClick={handleSave}
          >
            💾 บันทึก
          </button>
        </div>
      </div>
    </div>
  );
}
