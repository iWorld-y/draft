import React, { useState, useEffect } from 'react';
import Dashboard from './pages/Dashboard';
import DictionaryUpload from './pages/DictionaryUpload';
import Learning from './pages/Learning';
import './styles/global.css';

type Page = 'dashboard' | 'upload' | 'learn';

const App: React.FC = () => {
  const [currentPage, setCurrentPage] = useState<Page>('dashboard');
  const [dictId, setDictId] = useState<number>(1);

  useEffect(() => {
    const path = window.location.pathname;
    const searchParams = new URLSearchParams(window.location.search);
    
    if (path === '/upload') {
      setCurrentPage('upload');
    } else if (path === '/learn') {
      setCurrentPage('learn');
      const id = searchParams.get('dictId');
      if (id) {
        setDictId(parseInt(id, 10));
      }
    } else {
      setCurrentPage('dashboard');
    }
  }, []);

  const navigateTo = (page: Page, params?: { dictId?: number }) => {
    setCurrentPage(page);
    if (params?.dictId) {
      setDictId(params.dictId);
    }
    
    // Update URL
    let url = '/';
    if (page === 'upload') {
      url = '/upload';
    } else if (page === 'learn') {
      url = params?.dictId ? `/learn?dictId=${params.dictId}` : '/learn';
    }
    window.history.pushState({}, '', url);
  };

  // Override window.location for navigation
  useEffect(() => {
    const originalHref = window.location.href;
    
    const handleClick = (e: MouseEvent) => {
      const target = e.target as HTMLElement;
      const anchor = target.closest('a');
      
      if (anchor) {
        const href = anchor.getAttribute('href');
        if (href && href.startsWith('/')) {
          e.preventDefault();
          if (href === '/upload') {
            navigateTo('upload');
          } else if (href.startsWith('/learn')) {
            const url = new URL(href, window.location.origin);
            const dictId = url.searchParams.get('dictId');
            navigateTo('learn', { dictId: dictId ? parseInt(dictId, 10) : 1 });
          } else {
            navigateTo('dashboard');
          }
        }
      }
    };

    document.addEventListener('click', handleClick);
    return () => document.removeEventListener('click', handleClick);
  }, []);

  const renderPage = () => {
    switch (currentPage) {
      case 'upload':
        return <DictionaryUpload />;
      case 'learn':
        return <Learning dictId={dictId} />;
      case 'dashboard':
      default:
        return <Dashboard />;
    }
  };

  return (
    <div className="app">
      {renderPage()}
    </div>
  );
};

export default App;
