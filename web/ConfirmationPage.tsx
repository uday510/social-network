import { useParams, useNavigate } from "react-router-dom";
import { useState } from "react";
import { API_URL } from "./src/App.tsx";

export const ConfirmationPage = () => {
    const { token = "" } = useParams();
    const navigate = useNavigate();
    const [loading, setLoading] = useState(false);
    const [isActive, setIsActive] = useState(false);

    const handleConfirm = async () => {
        setLoading(true);
        try {
            const res = await fetch(`${API_URL}/users/activate/${token}`, {
                method: "PUT",
            });

            if (!res.ok) {
                alert("Failed to confirm email.");
                return;
            }

            navigate("/");
        } catch (err) {
            alert("Something went wrong.");
        } finally {
            setLoading(false);
        }
    };

    return (
        <div style={styles.container}>
            <button
                style={{
                    ...styles.button,
                    backgroundColor: isActive ? "#ffffff" : "#333",
                    color: isActive ? "#333" : "#ffffff",
                    opacity: loading ? 0.6 : 1,
                    cursor: loading ? "not-allowed" : "pointer",
                }}
                onClick={handleConfirm}
                disabled={loading}
                onMouseDown={() => setIsActive(true)}
                onMouseUp={() => setIsActive(false)}
                onMouseLeave={() => setIsActive(false)}
            >
                {loading ? "Confirming..." : "Confirm Email"}
            </button>
        </div>
    );
};

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
    button: {
        border: "2px solid #333",
        padding: "14px 28px",
        borderRadius: "6px",
        fontSize: "1rem",
        transition: "all 0.2s ease",
    },
};

export default ConfirmationPage;