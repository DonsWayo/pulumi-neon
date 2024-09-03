using System.Collections.Generic;
using System.Linq;
using Pulumi;
using neon = Pulumi.neon;

return await Deployment.RunAsync(() => 
{
    var myRandomResource = new neon.Random("myRandomResource", new()
    {
        Length = 24,
    });

    return new Dictionary<string, object?>
    {
        ["output"] = 
        {
            { "value", myRandomResource.Result },
        },
    };
});

