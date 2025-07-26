import './App.css';

export const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/v1';

function App() {
    return (
        <div style={styles.container}>
            <div style={styles.content}>SocialNetwork</div>
        </div>
    );
}

const styles: { [key: string]: React.CSSProperties } = {
    container: {
        position: "fixed",
        top: 0,
        left: 0,
        width: "100vw",
        height: "100vh",
        backgroundColor: "#ffffff",
        display: "flex",
        justifyContent: "center",
        alignItems: "center",
        zIndex: 9999,
    },
    content: {
        fontSize: '1.5rem',
        color: '#333',
    },
};

export default App;