// Pulumi Command Provider .NET SDK
// Copyright 2020, Mitchell Maler.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

using System;
using System.IO;
using System.Reflection;
using Pulumi;

namespace Pulumi.Command
{
    static class Utilities
    {
        public static string? GetEnv(params string[] names)
        {
            foreach (var n in names)
            {
                var value = Environment.GetEnvironmentVariable(n);
                if (value != null)
                {
                    return value;
                }
            }
            return null;
        }

        static string[] trueValues = { "1", "t", "T", "true", "TRUE", "True" };
        static string[] falseValues = { "0", "f", "F", "false", "FALSE", "False" };
        public static bool? GetEnvBoolean(params string[] names)
        {
            var s = GetEnv(names);
            if (s != null)
            {
                if (Array.IndexOf(trueValues, s) != -1)
                {
                    return true;
                }
                if (Array.IndexOf(falseValues, s) != -1)
                {
                    return false;
                }
            }
            return null;
        }

        public static int? GetEnvInt32(params string[] names) => int.TryParse(GetEnv(names), out int v) ? (int?)v : null;

        public static double? GetEnvDouble(params string[] names) => double.TryParse(GetEnv(names), out double v) ? (double?)v : null;

        public static InvokeOptions WithVersion(this InvokeOptions? options)
        {
            if (options?.Version != null)
            {
                return options;
            }
            return new InvokeOptions
            {
                Parent = options?.Parent,
                Provider = options?.Provider,
                Version = Version,
            };
        }

        private readonly static string version;
        public static string Version => version;

        static Utilities()
        {
            var assembly = typeof(Utilities).GetTypeInfo().Assembly;
            using var stream = assembly.GetManifestResourceStream("Pulumi.Command.version.txt");
            using var reader = new StreamReader(stream ?? throw new NotSupportedException("Missing embedded version.txt file"));
            version = reader.ReadToEnd().Trim();
        }
    }
}
