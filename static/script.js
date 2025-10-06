async function divine() {
    const question = document.getElementById('question').value.trim();
    const divineBtn = document.getElementById('divineBtn');
    const loading = document.getElementById('loading');
    const result = document.getElementById('result');
    const error = document.getElementById('error');

    error.style.display = 'none';
    result.style.display = 'none';
    loading.style.display = 'block';
    divineBtn.disabled = true;

    try {
        const response = await fetch('/api/divine', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                question: question || '我今天的运势如何？',
                reader: 'deepseek'
            }),
        });

        const data = await response.json();

        if (data.error) {
            throw new Error(data.error);
        }

        const cardsDiv = document.getElementById('cards');
        cardsDiv.innerHTML = '';
        data.cards.forEach(card => {
            const cardDiv = document.createElement('div');
            cardDiv.className = 'card';
            cardDiv.textContent = `${card.name_zh}（${card.position === 'Upright' ? '正位' : '逆位'}）`;
            cardsDiv.appendChild(cardDiv);
        });

        document.getElementById('reading').textContent = data.reading;

        const img = document.getElementById('resultImage');
        img.src = 'data:image/jpeg;base64,' + data.image;

        loading.style.display = 'none';
        result.style.display = 'block';

    } catch (err) {
        loading.style.display = 'none';
        error.style.display = 'block';
        error.textContent = '占卜失败：' + err.message;
    } finally {
        divineBtn.disabled = false;
    }
}

function reset() {
    document.getElementById('result').style.display = 'none';
    document.getElementById('question').value = '';
    document.getElementById('question').focus();
}

document.getElementById('question').addEventListener('keydown', (e) => {
    if (e.key === 'Enter' && e.ctrlKey) {
        divine();
    }
});
