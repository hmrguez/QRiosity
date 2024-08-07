import {gql, useQuery} from "@apollo/client";
import {useState} from "react";
import {useNavigate} from "react-router-dom";

const LOGIN_QUERY = gql`
    query Login($username: String!, $password: String!) {
        login(username: $username, password: $password) {
            token
            user {
                id
                name
            }
        }
    }
`;

const Login = () => {
    const [username, setUsername] = useState('');
    const [password, setPassword] = useState('');
    const {data, loading, error, refetch} = useQuery(LOGIN_QUERY, {
        variables: {username, password},
        skip: true,
    });
    const navigate = useNavigate();

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        const response = await refetch();
        if (response.data) {
            localStorage.setItem('token', response.data.login.token);
            navigate('/');
        }
    };

    return (
        <form onSubmit={handleSubmit}>
            <div>
                <label>Username</label>
                <input type="text" value={username} onChange={(e) => setUsername(e.target.value)} />
            </div>
            <div>
                <label>Password</label>
                <input type="password" value={password} onChange={(e) => setPassword(e.target.value)} />
            </div>
            <button type="submit" disabled={loading}>Login</button>
            {error && <p>Error: {error.message}</p>}
        </form>
    );
};

export default Login;