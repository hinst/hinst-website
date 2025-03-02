import { useLocation } from 'react-router';
import icon from './icon.webp';

export default function Header(props: {title: string}) {
    const location = useLocation();
    return <div style={{display: 'flex', alignItems: 'center', gap: 10}}>
        <img src={icon} width={56} height={56} alt='icon'/>
        <h6>{props.title}</h6>
    </div>;
}