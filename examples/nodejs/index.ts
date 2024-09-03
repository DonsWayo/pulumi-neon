import * as pulumi from "@pulumi/pulumi";
import * as neon from "@pulumi/neon";

const myRandomResource = new neon.Random("myRandomResource", {length: 24});
export const output = {
    value: myRandomResource.result,
};
