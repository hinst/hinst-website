import { NavLink, HashRouter } from 'react-router';
import Header from './header';
import { useState } from 'react';

const PAGE_TITLE = 'Hidden Personal Website';

export default function App() {
    const pageTitle = useState(PAGE_TITLE);
    return <HashRouter>
        <Header/>
        <NavLink to='/personal-goals' className='ms-btn ms-outline ms-primary'>My Personal Goals</NavLink>
    </HashRouter>;
}
