import * as pulumi from "@pulumi/pulumi";
import * as neon from "@pulumi/neon";

// Get the Neon API key from Pulumi config
const config = new pulumi.Config();
//const neonApiKey = config.requireSecret("neonApiKey");

// Configure the Neon provider with the API key
const neonProvider = new neon.Provider("neonProvider", {
    apiKey: "c1v9gfm4f2ug1hsacotgaix1ibvrbw57yd0f34msbdnja2zy25tl8qbv52q71w1f",
});

const project = new neon.Project("MyNeonProject", {
    name: "MyNeonProject",
    regionId: "aws-us-west-2",
}, { provider: neonProvider });

const branch = new neon.Branch("myBranch", {
    projectId: project.id,
    name: "main",
}, { provider: neonProvider });

const database = new neon.Database("myDatabase", {
    projectId: project.id,
    branchId: branch.id,
    name: "mydb",
}, { provider: neonProvider });

export const neonProjectId = project.id;
export const neonBranchId = branch.id;
export const neonDatabaseName = database.name;