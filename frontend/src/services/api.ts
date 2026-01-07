const API_BASE_URL: string = process.env.REACT_APP_API_URL || '/api';

interface ApiData {
  [key: string]: any;
}

export const api = {
  async getData(): Promise<ApiData> {
    const response: Response = await fetch(`${API_BASE_URL}/data`);
    if (!response.ok) throw new Error('Failed to fetch data');
    return response.json();
  },

  async postData(data: ApiData): Promise<ApiData> {
    const response: Response = await fetch(`${API_BASE_URL}/data`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(data),
    });
    if (!response.ok) throw new Error('Failed to post data');
    return response.json();
  }
};