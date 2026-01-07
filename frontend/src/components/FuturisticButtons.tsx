import { useState, useEffect } from 'react';
import './FuturisticButtons.css';

function FuturisticButtons() {
  const [response, setResponse] = useState<string>('');
  const [loading, setLoading] = useState<boolean>(false);
  const [fadeOut, setFadeOut] = useState<boolean>(false);

  const fetchData = async (endpoint: string) => {
    setLoading(true);
    setFadeOut(false);
    try {
      const res = await fetch(`/api${endpoint}`, {
        cache: 'no-cache',
        headers: {
          'Cache-Control': 'no-cache'
        }
      });
      const text = await res.text();
      setResponse(text);
    } catch (error) {
      setResponse(`Error: ${error}`);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    if (response && !loading) {
      const fadeTimer = setTimeout(() => {
        setFadeOut(true);
      }, 800);

      const clearTimer = setTimeout(() => {
        setResponse('');
        setFadeOut(false);
      }, 2000);

      return () => {
        clearTimeout(fadeTimer);
        clearTimeout(clearTimer);
      };
    }
  }, [response, loading]);

  const handleClickButton1 = () => {
    fetchData('/');
  };

  const handleClickButton2 = () => {
    fetchData('/hello');
  };

  const handleClickButton3 = () => {
    fetchData('/bye');
  };

  return (
    <div className="button-container">
      <div className="button-box">
        <button className="futuristic-btn" onClick={handleClickButton2}>Hello</button>
        <button className="futuristic-btn" onClick={handleClickButton3}>Bye</button>
        <button className="futuristic-btn" onClick={handleClickButton1}>Knock</button>
      </div>

      {(response || loading) && (
        <div className={`response-box ${fadeOut ? 'fade-out' : ''}`}>
          {loading ? 'Loading...' : response}
        </div>
      )}
    </div>
  );
}

export default FuturisticButtons;
