import { useLocation } from 'react-router';

export default function Header() {
    const location = useLocation();
    return <h1>Hidden Personal Website</h1>;
}