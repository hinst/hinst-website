import { useLocation } from 'react-router';

export default function Header(props: {title: string}) {
    const location = useLocation();
    return <h6>{props.title}</h6>;
}