import { NavLink, HashRouter } from 'react-router';
import Header from './header';
import { useEffect, useState } from 'react';

const PAGE_TITLE = 'Showcase Website';

export default function App() {
    const [pageTitle, setPageTitle] = useState(PAGE_TITLE);
    useEffect(() => {
        document.title = pageTitle;
    }, [pageTitle]);
    return <HashRouter>
        <div style={{marginBottom: 10}}>
            <Header title={pageTitle}/>
        </div>
        <NavLink to='/personal-goals' className='ms-btn ms-outline ms-primary'>My Personal Goals</NavLink>
    </HashRouter>;
}
