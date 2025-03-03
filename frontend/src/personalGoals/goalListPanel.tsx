import { useEffect, useState } from 'react';
import { GoalHeader } from './goalHeader';
import { API_URL } from '../api';

export default function GoalList() {
    const [goals, setGoals] = useState(new Array<GoalHeader>());
    async function loadGoals() {
        const response = await fetch(API_URL + '/goals');
        if (!response.ok)
            throw new Error(response.statusText);
        const data = await response.json();
        setGoals(data);
    };
    useEffect(
        () => { loadGoals(); },
        [],
    );
    return <div>
    </div>;
}