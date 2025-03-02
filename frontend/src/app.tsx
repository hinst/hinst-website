import { NavLink, HashRouter } from 'react-router';

export default function App() {
    return <HashRouter>
        <NavLink to='/personal-goals' className='ms-btn'>Personal Goals</NavLink>
    </HashRouter>;
}
