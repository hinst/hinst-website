import { createRoot } from 'react-dom/client';
import 'minstyle.io/dist/css/minstyle.io.css';
import './index.css';

const appElement = document.getElementById('app');
if (!appElement)
    throw new Error('Cannot find element app');

const root = createRoot(appElement);
root.render(<h1>Hello, world</h1>);