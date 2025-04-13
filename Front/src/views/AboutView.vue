<template>
  <div class="about-page">
    <div class="about-hero">
      <div class="container">
        <h1 class="about-title">À PROPOS</h1>
      </div>
    </div>

    <div class="container about-content">
      <div class="about-section">
        <div class="about-text">
          <h2 class="section-title">NOTRE HISTOIRE</h2>
          <p>
            BRUT DESIGN a été fondé en 2023 avec une vision claire : créer des
            meubles et accessoires de design néobrutaliste qui allient
            esthétique brute et fonctionnalité. Inspirés par les mouvements
            architecturaux brutalistes des années 50-70, nous créons des pièces
            qui célèbrent la beauté des matériaux dans leur forme la plus
            authentique.
          </p>
          <p>
            Notre approche du design est ancrée dans l'honnêteté des matériaux,
            la géométrie audacieuse et les textures expressives. Chaque pièce de
            notre collection raconte une histoire de créativité sans compromis
            et d'artisanat méticuleux.
          </p>
        </div>

        <div class="about-image image-box-1">
          <div class="image-container">
            <img
              src="https://images.unsplash.com/photo-1618220179428-22790b461013?ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D&auto=format&fit=crop&w=1000&q=80"
              alt="Atelier de design"
            />
          </div>
        </div>
      </div>

      <div class="about-section reverse">
        <div class="about-text">
          <h2 class="section-title">NOTRE PHILOSOPHIE</h2>
          <p>
            Le néobrutalisme est plus qu'un style pour nous, c'est une
            philosophie. Nous croyons en la création de pièces qui défient les
            tendances éphémères et qui perdurent, tant dans leur durabilité
            physique que dans leur pertinence esthétique.
          </p>
          <p>
            Nous rejetons la production de masse et l'obsolescence programmée. À
            la place, nous proposons des objets qui sont conçus pour durer et
            vieillir avec caractère. Chaque imperfection et variation est
            célébrée comme une marque d'authenticité, rappelant que derrière
            chaque création se trouve un processus humain.
          </p>
        </div>

        <div class="about-image image-box-2">
          <div class="image-container">
            <img
              src="https://images.unsplash.com/photo-1586023492125-27b2c045efd7?ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D&auto=format&fit=crop&w=1000&q=80"
              alt="Design brutaliste"
            />
          </div>
        </div>
      </div>

      <div class="about-values">
        <h2 class="section-title centered">NOS VALEURS</h2>

        <div class="values-grid">
          <div class="value-card">
            <h3>AUTHENTICITÉ</h3>
            <p>
              Nous valorisons l'honnêteté des matériaux et des processus, sans
              masquer leur nature intrinsèque.
            </p>
          </div>

          <div class="value-card">
            <h3>DURABILITÉ</h3>
            <p>
              Nos créations sont conçues pour durer, tant dans leur construction
              que dans leur design intemporel.
            </p>
          </div>

          <div class="value-card">
            <h3>EXPRESSION</h3>
            <p>
              Nous encourageons l'audace et l'expressivité dans le design,
              refusant les conventions esthétiques dominantes.
            </p>
          </div>

          <div class="value-card">
            <h3>ARTISANAT</h3>
            <p>
              Chaque pièce est réalisée avec soin, en valorisant le savoir-faire
              et l'attention aux détails.
            </p>
          </div>
        </div>
      </div>

      <div class="contact-section">
        <h2 class="section-title">CONTACTEZ-NOUS</h2>
        <div class="contact-info">
          <p><strong>Email:</strong> contact@brutdesign.fr</p>
          <p><strong>Téléphone:</strong> +33 1 23 45 67 89</p>
          <p>
            <strong>Showroom:</strong> 42 Rue de l'Innovation, 75003 Paris,
            France
          </p>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.about-page {
  margin-bottom: 2rem;
}

.about-hero {
  background-color: #000;
  color: #fff;
  padding: 5rem 0;
  position: relative;
  overflow: hidden;
}

.about-hero::before {
  content: "";
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: repeating-linear-gradient(
    45deg,
    rgba(255, 255, 255, 0.05),
    rgba(255, 255, 255, 0.05) 10px,
    transparent 10px,
    transparent 20px
  );
}

.about-title {
  font-size: 4rem;
  font-weight: 900;
  letter-spacing: -2px;
  margin: 0;
  transform: skew(-5deg);
  display: inline-block;
  padding: 10px 20px;
  border: 4px solid #fff;
  position: relative;
  z-index: 2;
}

.about-content {
  padding: 3rem 0;
}

.about-section {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 3rem;
  margin-bottom: 5rem;
  align-items: center;
}

.about-section.reverse {
  grid-template-columns: 1fr 1fr;
  direction: rtl;
}

.about-section.reverse > * {
  direction: ltr;
}

.about-text {
  padding: 2rem;
  background-color: #f5f5f5;
  border: 4px solid #000;
  box-shadow: 8px 8px 0 #000;
  transform: rotate(-1deg);
}

.section-title {
  font-size: 2.5rem;
  margin-bottom: 1.5rem;
  font-weight: 900;
  letter-spacing: -1px;
  text-transform: uppercase;
  position: relative;
  display: inline-block;
}

.section-title::after {
  content: "";
  position: absolute;
  bottom: -5px;
  left: 0;
  width: 100%;
  height: 3px;
  background-color: #000;
}

.section-title.centered {
  text-align: center;
  display: block;
}

.section-title.centered::after {
  left: 50%;
  transform: translateX(-50%);
  width: 150px;
}

.about-text p {
  margin-bottom: 1rem;
  font-size: 1.1rem;
  line-height: 1.7;
}

.image-container {
  overflow: hidden;
  border: 4px solid #000;
  height: 100%;
}

.image-container img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  display: block;
  transition: transform 0.5s;
}

.image-container:hover img {
  transform: scale(1.03);
}

.image-box-1 {
  transform: rotate(2deg);
  box-shadow: -8px 8px 0 #000;
  height: 400px;
}

.image-box-2 {
  transform: rotate(-2deg);
  box-shadow: 8px 8px 0 #000;
  height: 400px;
}

.about-values {
  margin: 5rem 0;
  padding: 2rem 0;
}

.values-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
  gap: 2rem;
  margin-top: 3rem;
}

.value-card {
  background-color: #f5f5f5;
  border: 4px solid #000;
  padding: 1.5rem;
  transform: skew(0deg) rotate(0deg);
  transition: transform 0.3s, box-shadow 0.3s;
  box-shadow: 5px 5px 0 #000;
}

.value-card:nth-child(odd) {
  transform: rotate(-1deg);
}

.value-card:nth-child(even) {
  transform: rotate(1deg);
}

.value-card:hover {
  transform: translateY(-5px);
  box-shadow: 8px 8px 0 #000;
}

.value-card h3 {
  font-size: 1.4rem;
  margin-bottom: 1rem;
  font-weight: 800;
  letter-spacing: -0.5px;
}

.contact-section {
  margin-top: 4rem;
  padding: 2rem;
  background-color: #f5f5f5;
  border: 4px solid #000;
  transform: rotate(-1deg);
  box-shadow: 8px 8px 0 #000;
}

.contact-info {
  margin-top: 2rem;
  font-size: 1.1rem;
}

.contact-info p {
  margin-bottom: 1rem;
}

@media (max-width: 900px) {
  .about-section,
  .about-section.reverse {
    grid-template-columns: 1fr;
    gap: 1.5rem;
  }

  .image-box-1,
  .image-box-2 {
    height: 300px;
  }

  .about-title {
    font-size: 3rem;
  }
}
</style>
