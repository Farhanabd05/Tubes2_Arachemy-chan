import json

# Baca JSON asli dari file (misal 'data_input.json')
with open('recipes.json', 'r', encoding='utf-8') as f:
    data = json.load(f)

# Transformasi: buat list baru dengan kunci diurutkan
transformed = [
    {
        "input": item["input"],
        "output": item["output"]
    }
    for item in data
]

# Simpan hasil ke file baru (misal 'data_output.json')
with open('data_output.json', 'w', encoding='utf-8') as f:
    json.dump(transformed, f, ensure_ascii=False, indent=2)