import { useLocation } from 'react-router';

export default function Header() {
    const location = useLocation();
    return <h6>Hidden Personal Website</h6>;
}