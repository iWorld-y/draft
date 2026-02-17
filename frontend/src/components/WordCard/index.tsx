import React, { useState } from 'react';
import type { Word } from '../../services/learning';
import './WordCard.css';

interface WordCardProps {
  word: Word;
  onReveal: () => void;
}

const WordCard: React.FC<WordCardProps> = ({ word, onReveal }) => {
  const [revealed, setRevealed] = useState(false);
  
  const handleReveal = () => {
    setRevealed(true);
    onReveal();
  };
  
  return (
    <div className="word-card">
      <div className="word-header">
        <h1 className="word-text">{word.word}</h1>
        {word.phonetic && <span className="phonetic">{word.phonetic}</span>}
      </div>
      
      {!revealed ? (
        <button 
          className="reveal-button"
          onClick={handleReveal}
        >
          显示释义
        </button>
      ) : (
        <div className="word-meaning">
          {word.meaning.definitions.map((def, idx) => (
            <div key={idx} className="definition-item">
              <span className="pos">{def.pos}</span>
              <span className="text">{def.text}</span>
            </div>
          ))}
          
          {word.example && (
            <div className="example">
              <label>例句：</label>
              <p>{word.example}</p>
            </div>
          )}
        </div>
      )}
    </div>
  );
};

export default WordCard;
