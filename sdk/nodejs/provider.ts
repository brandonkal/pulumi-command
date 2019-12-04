import * as pulumi from '@pulumi/pulumi'

/**
 * Provides an OpenFaaS Provider resource.
 */
export class Provider extends pulumi.ProviderResource {
  /**
   * Create a provider resource with the given unique name, arguments, and options.
   *
   * @param name The _unique_ name of the resource.
   * @param args The arguments to use to populate this provider's configuration.
   * @param opts A bag of options that control this resource's behavior.
   */
  constructor(name: string, args: ProviderArgs, opts?: pulumi.ResourceOptions) {
    super(
      'openfaas',
      name,
      {
        endpoint: args.endpoint,
        username: args.username,
        password: args.password,
        tlsSkipVerify: args.tlsSkipVerify,
      },
      opts
    )
  }
}

/**
 * The set of arguments for constructing a Provider resource.
 */
export interface ProviderArgs {
  readonly endpoint: pulumi.Input<string>
  readonly username?: pulumi.Input<string>
  readonly password?: pulumi.Input<string>
  readonly tlsSkipVerify?: pulumi.Input<boolean>
}
