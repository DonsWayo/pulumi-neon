import pulumi
import pulumi_neon as neon

my_random_resource = neon.Random("myRandomResource", length=24)
pulumi.export("output", {
    "value": my_random_resource.result,
})
