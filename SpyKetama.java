import java.net.InetSocketAddress;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.List;
import java.util.UUID;
import net.spy.memcached.ConnectionFactory;
import net.spy.memcached.DefaultHashAlgorithm;
import net.spy.memcached.KetamaConnectionFactory;
import net.spy.memcached.KetamaNodeLocator;
import net.spy.memcached.MemcachedNode;
import net.spy.memcached.util.DefaultKetamaNodeLocatorConfiguration;

public class SpyKetama {
	public static void main(String []args) {
		ConnectionFactory cf = new KetamaConnectionFactory();

		DefaultKetamaNodeLocatorConfiguration config =
			new DefaultKetamaNodeLocatorConfiguration();

		List<InetSocketAddress> addrs = Arrays.asList(
			new InetSocketAddress("127.0.0.1", 11211),
			new InetSocketAddress("127.0.0.1", 11212),
			new InetSocketAddress("127.0.0.1", 11213),
			new InetSocketAddress("127.0.0.1", 11214),
			new InetSocketAddress("localhost", 11211),
			new InetSocketAddress("localhost", 11212),
			new InetSocketAddress("localhost", 11213),
			new InetSocketAddress("localhost", 11214));

		List<MemcachedNode> nodes = new ArrayList<MemcachedNode>();
		for (InetSocketAddress a : addrs) {
			nodes.add(cf.createMemcachedNode(a, null, 1024));
		}

		KetamaNodeLocator l = new KetamaNodeLocator(
			nodes,
			DefaultHashAlgorithm.KETAMA_HASH,
			config);

		// Sprinkle a bit of magic on this number...
		int magicReps = config.getNodeRepetitions() / 4;

		System.out.println("package gomcketama\n");
		System.out.println("var (");

		System.out.println("	nodeKeys = map[string][]string{");
		for (MemcachedNode n : nodes) {
			InetSocketAddress a = (InetSocketAddress)n.getSocketAddress();

			System.out.format("		\"%s:%d\": []string{\n",
				a.getHostString(), a.getPort());

			for (int i = 0; i < magicReps; i++) {
				System.out.format("		\"%s\",\n", config.getKeyForNode(n, i));
			}

			System.out.println("	},");
		}
		System.out.println("	}");

		System.out.println("	kvToNode = map[string]string{");
		for (int i = 0; i < 20000; i++) {
			InetSocketAddress a = (InetSocketAddress)l
				.getPrimary(Integer.toString(i))
				.getSocketAddress();

			System.out.format("		\"%d\": \"%s:%d\",\n",
				i,
				a.getHostString(), a.getPort());
		}
		for (int i = 0; i < 2000; i++) {
			String uuid = UUID.randomUUID().toString();
			InetSocketAddress a = (InetSocketAddress)l
				.getPrimary(uuid)
				.getSocketAddress();

			System.out.format("		\"%s\": \"%s:%d\",\n",
				uuid,
				a.getHostString(), a.getPort());
		}
		System.out.println("	}");
		System.out.println(")");

		System.exit(0);
	}
}
