import React from 'react';
import './QualityButtons.css';

interface Quality {
  value: number;
  label: string;
  color: string;
}

const QUALITIES: Quality[] = [
  { value: 0, label: '完全不认识', color: '#ff4d4f' },
  { value: 1, label: '有印象', color: '#ff7a45' },
  { value: 2, label: '想起来了', color: '#ffa940' },
  { value: 3, label: '有些犹豫', color: '#ffc53d' },
  { value: 4, label: '轻松想起', color: '#95de64' },
  { value: 5, label: '脱口而出', color: '#52c41a' },
];

interface QualityButtonsProps {
  onSelect: (value: number) => void;
  disabled?: boolean;
}

const QualityButtons: React.FC<QualityButtonsProps> = ({ onSelect, disabled }) => {
  return (
    <div className="quality-buttons">
      <p className="hint">你对这个单词的掌握程度？</p>
      <div className="buttons-grid">
        {QUALITIES.map(q => (
          <button
            key={q.value}
            className="quality-btn"
            style={{ borderColor: q.color, '--hover-bg': q.color } as React.CSSProperties}
            onClick={() => onSelect(q.value)}
            disabled={disabled}
          >
            <span className="value" style={{ color: q.color }}>{q.value}</span>
            <span className="label">{q.label}</span>
          </button>
        ))}
      </div>
    </div>
  );
};

export default QualityButtons;
