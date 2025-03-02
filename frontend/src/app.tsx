import { NavLink, HashRouter } from 'react-router';
import Header from './header';

export default function App() {
    return <HashRouter>
        <Header/>
        <NavLink to='/personal-goals' className='ms-btn'>Personal Goals</NavLink>
    </HashRouter>;
}
