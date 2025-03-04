import { createRoot } from 'react-dom/client';
import 'minstyle.io/dist/css/minstyle.io.css';
import './index.css';
import App from './app';

const appElement = document.getElementById('app');
if (!appElement)
	throw new Error('Cannot find element app');

const root = createRoot(appElement);
root.render(<App />);