import React, { useEffect, useRef } from 'react';

interface User {
  name: string;
  imageUrl: string;
  imageSize: number;
}

const user: User = {
  name: 'Hedy Lamarr',
  imageUrl: 'https://i.imgur.com/yXOvdOSs.jpg',
  imageSize: 90,
};

function Lissajous(): React.JSX.Element {
  const canvasRef = useRef<HTMLCanvasElement>(null);

  useEffect(() => {
    const canvas = canvasRef.current;
    if (!canvas) return;
    const ctx = canvas.getContext('2d');
    if (!ctx) return;

    const cycles = 5;
    const res = 0.001;
    const size = 100;
    const nframes = 64;
    const delay = 80; // 8 * 10ms
    const freq = Math.random() * 3.0;
    let phase = 0.0;
    let frame = 0;
    let lastTime = 0;
    let requestId: number;

    const render = (time: number) => {
      if (time - lastTime >= delay) {
        lastTime = time;

        // Clear canvas
        ctx.fillStyle = 'white';
        ctx.fillRect(0, 0, 2 * size + 1, 2 * size + 1);

        // Draw Lissajous
        ctx.fillStyle = 'black';
        for (let t = 0.0; t < cycles * 2 * Math.PI; t += res) {
          const x = Math.sin(t);
          const y = Math.sin(t * freq + phase);
          // Set pixel (using fillRect as a proxy for SetColorIndex)
          ctx.fillRect(size + x * size, size + y * size, 1, 1);
        }

        phase += 0.1;
        frame = (frame + 1) % nframes;
      }
      requestId = requestAnimationFrame(render);
    };

    requestId = requestAnimationFrame(render);
    return () => cancelAnimationFrame(requestId);
  }, []);

  return (
    <div style={{ marginTop: '20px' }}>
      <h3>Lissajous Animation</h3>
      <canvas
        ref={canvasRef}
        width={201}
        height={201}
        style={{ border: '1px solid #ccc' }}
      />
    </div>
  );
}

interface Article {
  title: string;
  content: string;
}

function ArticleList() {
  const [articles, setArticles] = React.useState<Article[]>([]);

  useEffect(() => {
    fetch('/articles')
      .then(res => res.json())
      .then(data => {
          if (data.articles) {
              setArticles(data.articles);
          }
      })
      .catch(err => console.error(err));
  }, []);

  return (
    <div style={{ marginTop: '20px', borderTop: '1px solid #ccc', paddingTop: '20px' }}>
      <h3>Article List</h3>
      {articles.length === 0 ? (
        <p>No articles found.</p>
      ) : (
        <ul style={{ listStyleType: 'none', padding: 0 }}>
          {articles.map((article, index) => (
            <li key={index} style={{ marginBottom: '15px', border: '1px solid #eee', padding: '10px' }}>
              <h4 style={{ margin: '0 0 5px 0' }}>{article.title}</h4>
              <pre style={{ whiteSpace: 'pre-wrap', backgroundColor: '#f5f5f5', padding: '5px' }}>{article.content}</pre>
            </li>
          ))}
        </ul>
      )}
    </div>
  );
}

export default function Profile(): React.JSX.Element {
  return (
    <>
      <h1>{user.name}</h1>
      <img
        className="avatar"
        src={user.imageUrl}
        alt={'Photo of ' + user.name}
        style={{
          width: user.imageSize,
          height: user.imageSize
        }}
      />
      <Lissajous />
      <ArticleList />
    </>
  );
}