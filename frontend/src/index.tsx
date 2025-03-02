import { createRoot } from 'react-dom/client';
import 'minstyle.io/dist/css/minstyle.io.css';

// Clear the existing HTML content
document.body.innerHTML = '<div id="app"></div>';

const appElement = document.getElementById('app');
if (!appElement)
    throw new Error('Cannot find element app');

const root = createRoot(appElement);
root.render(<h1>Hello, world</h1>);