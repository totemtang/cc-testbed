import sys

import numpy as np
from sklearn.svm import SVC
from sklearn import tree
from sklearn.ensemble import RandomForestClassifier, AdaBoostClassifier
from sklearn.multiclass import OneVsRestClassifier

START = 5
FEATURELEN = 6
PARTAVG = 0
PARTSKEW = 1
RECAVG = 2
LATENCY = 3
READRATE = 4
CONFRATE = 5

class SmallbankPart(object):

	def __init__(self, f):
		self.clf = self.train(f)

	# X is feature, while Y is label
	def train(self, f):
		X = []
		Y = []
		for line in open(f):
			columns = [float(x) for x in line.strip().split('\t')[START:]]
			tmpX = []
			tmpX.extend(columns[:FEATURELEN])
			#tmpX.extend(columns[3:7])
			X.append(tmpX)
			if (columns[FEATURELEN] == 0):
				Y.extend([0])
			else:
				Y.extend([1])
		clf = tree.DecisionTreeClassifier(max_depth=6)
		clf = clf.fit(X, Y)
		return clf

	def Predict(self, partAvg, partSkew, recAvg, hitRate, readRate, confRate):
		X = [[partAvg, partSkew, recAvg, hitRate, readRate, confRate]]
		return self.clf.predict(X)[0]

class SmallbankOCC(object):

	def __init__(self, f):
		self.clf = self.train(f)

	def train(self, f):
		X = []
		Y = []
		for line in open(f):
			columns = [float(x) for x in line.strip().split('\t')[START:]]
			tmpY = columns[FEATURELEN:]
			if (columns[FEATURELEN] != 0):
				if (len(columns) <= FEATURELEN + 1 or (len(columns) > FEATURELEN + 1 and columns[FEATURELEN+1] != 0)):
					tmp = []
					tmp.extend(columns[RECAVG:FEATURELEN])
					X.append(tmp)
					Z = [0, 0]
					for y in tmpY:
						Z[int(y) - 1] = 1
					Y.append(Z)
			#if (len(columns) > FEATURELEN + 1):
			#	Y.extend([3])
			#elif (columns[FEATURELEN] == 1):
			#	Y.extend([1])
			#elif (columns[FEATURELEN] == 2):
			#	Y.extend([2])
		#clf = tree.DecisionTreeClassifier(max_depth=6)
		clf = OneVsRestClassifier(tree.DecisionTreeClassifier(max_depth=6))
		clf = clf.fit(np.array(X), np.array(Y))
		return clf

	def Predict(self, recAvg, hitRate, readRate, confRate):
		X = [[recAvg, hitRate, readRate, confRate]]
		Y = self.clf.predict(X)[0]
		if (Y[0] == 1 and Y[1] == 1):
			return 3
		elif (Y[1] == 0):
			return 1
		else:
			return 2

