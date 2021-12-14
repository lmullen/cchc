package main

import (
	"fmt"
	"log"
)

func main() {

	type Doc struct {
		lang string
		text string
	}

	es := Doc{lang: "es", text: "AL EXCMO. SEÑOR DON JOSÉ RAMON DEMETRIO FERNANDEZ Y MARTINEZ,,,Marqués de la Esperanza, Gran Cruz de la Real órden americana de Isabel la Católica, Comendador de número de la misma órden, ex-Diputado Constituyente, Teniente Coronel honorario del Batallon de Voluntarios de Puerto-Rico, Presidente del Centro hispano-ultramarino, etc., etc.,Excmo. Sr.:,Sabemos que al aparecer en público la,Historia de la Insurreccion de Lares,,un clamoreo se alzará por los que, haciendo poca justicia á los españoles nacidos en este suelo, creen que,puerto-riqueño,y,anti-español,son voces sinónimas; y se nos dirá que calumniamos al país, que somos enemigos de los hijos de esta provincia. Sea el nombre de V. E. escudo en que se emboten los tiros de nuestros detractores. Puerto-riqueño y ligado á su suelo natal por los lazos de la familia y de la propiedad, ninguno tiene mayores motivos ni mas derechò que V. E. para interesarse por el bienestar, por el verdadero progreso y prosperidad de Puerto-Rico. Nadie puede sostener que odiamos á nuestros hermanos de esta preciada Antilla, llevando á su frente nuestro libro el respetabilisimo nombre de un puerto-riqueño como V. E. Detestamos, sí, como V. E., como todos los leales de aquende y allende el mar á los hipócritas que, bajo especiosos pretestos políticos, buscan en las luchas de los partidos peninsulares los medios de debilitar aquí los elementos españoles y de dar vuelo por medio de exageradas reformas á los elementos contrarios á la nacionalidad y al órden; pero ni abrigamos sentimientos hostiles contra los puerto-riqueños que no reniegan de España, ni á nuestros mismos adversarios es nuestra mente calumniar en lo mas minimo.,,Otro motivo nos mueve á dedicar á V. E. nuestra modesta obra: V. E. es el jefe querido y patriótico del partido liberal-conservador, del partido español,sin condiciones;,y al rendir este dábil tributo de admiracion y gratitud á tan dignisimo patricio, creemos obsequiar á todos y á cada uno de los buenos españoles aquí nacidos ó residentes, porque V. E. representa y sintetiza las nobles aspiraciones de todos los amantes del órden y de la nacionalidad española en Puerto-Rico.,,Dignese V. E. aceptar esta nuestra humilde ofrenda y la distinguida consideracion con que somos de V. E. atentos y S. S. Q. B. S. M.,"}
	en := Doc{lang: "en", text: "The Lares Revolt of 1868 sought the abolition of slavery, freedom of the press and commerce, and the independence of Puerto Rico. Six hundred men, led by liberals, drew up a provisional constitution and declared the Puerto Rican Republic, but they were defeated in their first clash with Spanish troops. Despite the movements quick defeat, during the 20th century the revolt has come to be viewed as the beginning of Puerto Rico's struggle for independence. Pérez Morís opposed Puerto Rican independence, but his work has served as an important study of the revolt."}
	both := Doc{lang: "es+en", text: es.text + en.text}
	de := Doc{lang: "de", text: "s war ein wundervoller alter Glaube bei den Griechen, daß jedem neugeborenen Menschenwesen ein Stern am Himmel angezündet werde, der bei seinem Tod erlösche. Die Helligkeit und Größe des Gestirnes mochten der Bedeutung der Persönlichkeit entsprechen: so rühmte man vom König Mithradates, der drei Kriege gegen Rom geführt hat, bei seiner Geburt sei ein Komet erschienen, dessen Schweif den vierten Teil des Himmels überzog und siebzig Tage sichtbar blieb."}
	tricky := Doc{lang: "tricky", text: "This is in English. So it's not very long. Como estas? Estoy muy bien."}

	docs := []Doc{es, en, both, de, tricky}

	for _, d := range docs {
		_, w, err := CalculateLanguages(d.text)
		if err != nil {
			log.Println(err)
			break
		}
		fmt.Println("Actual language: ", d.lang, "Stats: ", w)

		fmt.Print("\n\n")

	}

}
