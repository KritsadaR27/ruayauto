// pages/api/keywords.js

export default function handler(req, res) {
  if (req.method === 'GET') {
    // Return sample data
    res.status(200).json([
      { keyword: "hello", response: "Hi there!" },
      { keyword: "help", response: "How can I assist you?" },
      { keyword: "ราคา", response: "ราคาสินค้าเริ่มต้นที่ 99 บาทค่ะ" },
      { keyword: "ส่งฟรีไหม", response: "เรามีบริการส่งฟรีเมื่อยอดครบ 500 บาทค่ะ" },
    ]);
  } else if (req.method === 'POST') {
    // Log incoming JSON
    console.log('Received keywords:', req.body);
    res.status(200).json({ ok: true });
  } else {
    res.setHeader('Allow', ['GET', 'POST']);
    res.status(405).end(`Method ${req.method} Not Allowed`);
  }
}
