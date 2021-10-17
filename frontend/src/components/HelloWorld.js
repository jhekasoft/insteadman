import React, { useState } from 'react';
import Modal from 'react-modal';

function HelloWorld() {
	const [showModal, setShowModal] = useState(false);
	const [result, setResult] = useState(null);
	const [games, setGames] = useState("");

	const handleOpenModal = () => {
		setShowModal(true);

		window.backend.basic().then((result) => setResult(result));
		window.backend.games().then((games) => setGames(games))
	};

	const handleCloseModal = () => {
		setShowModal(false);
		window.backend.games().then((games) => setGames(games))
	};

	return (
		<div className="App">
			<button onClick={() => handleOpenModal()} type="button">
				Hello
      </button>
			<Modal
				appElement={document.getElementById("app")}
				isOpen={showModal}
				contentLabel="Minimal Modal Example"
			>
				<p>{result}</p>
				<p>{games}</p>
				<button onClick={() => handleCloseModal()}>Close Modal</button>
			</Modal>
		</div>
	);
}

export default HelloWorld;
