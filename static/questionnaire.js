/*
var dag = {
  nodes: [],
  edges: [],
  questionnaire: {}
  given: []
};
*/

function displayQuestionnaire(questionnaire) {
  const graphContainer = d3.select('#questionnaireContainer');
  // Clear the container
  graphContainer.html('');
  const nodeTable = graphContainer.append('table').attr('class', 'table');

  // Create a title for the questionnaire
  graphContainer.append('h2').text('Questionnaire');
  // Create a form for the questions
  const questionnaireForm = graphContainer.append('ul').attr('class', 'list-group').attr('class', 'list-group-flush');
  questionnaire.ProbConds.forEach(
    (question, qIndex) => {
      // Create a list item for the question
      const questionListItem = questionnaireForm.append('li').attr('class', 'list-group-item');
      if (question.Conditions !== null || question.Negations !== null) {
        questionListItem.append('div').text(`What is the probability that ${question.InQuestion} given:`);
      } else {
        questionListItem.append('div').text(`What is the probability that ${question.InQuestion}?`);
      }

      if (question.Conditions !== null) {
        question.Conditions.forEach((text) => {
          questionListItem.append('div').text("- "+text);
        });
      }
      if (question.Negations !== null) {
        question.Negations.forEach((text) => {
          questionListItem.append('div').text("- "+text);
        });
      } // Create an input field for the integer answer
      const inputField = questionListItem.append('input').attr('type', 'number').attr('min', 0).attr('max', 100).on('change', (event) => {
        const answer = parseInt(inputField.node().value);
        dag.questionnaire.ProbConds[qIndex].Answer = answer;
      });
      inputField.node().value = question.Answer;
    },
  );
  // Add a submit button to the form
  questionnaireForm.append('br');
  questionnaireForm.append('div').append('button').attr('class', 'btn btn-primary').text('Calculate')
    .on('click', (e) => {
      e.preventDefault();
      // Prevent the default form submission behavior
      // d3.event.preventDefault();
      // Get the answers from the input fields
      const inputs = questionnaireForm.selectAll('input').nodes();
      questionnaire.ProbConds.forEach(
        (question, index) => {
          const answer = parseInt(inputs[index].value);
          dag.questionnaire.ProbConds[index].Answer = answer;
        },
      );
      // Submit the questionnaire with the answers to the server
      submitQuestionnaire(dag.questionnaire);
    });
}
function submitQuestionnaire(questionnaire) {
  const data = JSON.stringify({
    nodes: dag.nodes,
    questionnaire,
  });

  fetch(
    'questionnaire',
    {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: data,
    },
  ).then((response) => response.json()).then(
    (answers) => {
      alert(answers.map(answer => {
        return "Probability that "+answer.InQuestion+": "+answer.Result;
      }).join("\n"));
    },
  ).catch((error) => {
    console.error('Error:', error);
  });
}
