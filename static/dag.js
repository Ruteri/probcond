let dag = {
  nodes: [],
  edges: [],
  questionnaire: {},
  experiments: [],
  given: [],
};

window.dag = dag;

function addNode() {
  const nodeNameInput = document.getElementById('nodeName');
  const nodeNegationInput = document.getElementById('nodeNegation');
  if (nodeNameInput.value !== '') {
    const node = {
      value: nodeNameInput.value,
      negation: nodeNegationInput.value,
    };
    dag.nodes.push(node);
    nodeNameInput.value = '';
    nodeNegationInput.value = '';
  }
  updateGraph();
}
function addEdge() {
  const srcNodeSelect = document.getElementById('srcNode');
  const dstNodeSelect = document.getElementById('dstNode');
  const srcNodeId = parseInt(srcNodeSelect.value);
  const dstNodeId = parseInt(dstNodeSelect.value);
  if (srcNodeId !== dstNodeId) {
    const srcNode = dag.nodes[srcNodeId].value;
    const dstNode = dag.nodes[dstNodeId].value;
    const edge = {
      src: srcNode,
      dst: dstNode,
    };
    dag.edges.push(edge);
    srcNodeSelect.value = '';
    dstNodeSelect.value = '';
  }
  updateGraph();
}
function addExperiment(experiment, node) {
  dag.experiments.push([node, experiment]);
  updateGraph();
}
function exportData() {
  const data = JSON.stringify(dag);
  const blob = new Blob([data], {
    type: 'application/json',
  });
  const fileName = prompt('Enter file name:', 'dag.json');
  const link = document.createElement('a');
  link.href = URL.createObjectURL(blob);
  link.download = fileName;
  link.click();
}
function importData(event) {
  const file = event.target.files[0];
  const reader = new FileReader();
  reader.onload = function (e) {
    const contents = e.target.result;
    const importedDAG = JSON.parse(contents);
    // Update the existing DAG variable
    dag = importedDAG;
    if (dag.experiments === undefined) {
      dag.experiments = [];
    }
    if (dag.questionnaire.Experiments === undefined) {
      dag.questionnaire.Experiments = [];
    }
    // Add any necessary code to update the UI or perform other operations with the imported data
    updateGraph();
    displayQuestionnaire(dag.questionnaire);
  };
  reader.readAsText(file);
}
